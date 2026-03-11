#!/usr/bin/env python3
"""
Pulumi Neo Task Manager

Create and manage Pulumi Neo agent tasks via the REST API.
Requires PULUMI_ACCESS_TOKEN environment variable.
"""

import argparse
import json
import os
import subprocess
import sys
import time
from datetime import datetime, timezone
from typing import Optional

try:
    import requests
except ImportError:
    print("Error: requests library required. Install with: pip install requests")
    sys.exit(1)


BASE_URL = "https://api.pulumi.com/api/preview/agents"


def get_headers() -> dict:
    """Get required API headers with authentication."""
    token = os.environ.get("PULUMI_ACCESS_TOKEN")
    if not token:
        print("Error: PULUMI_ACCESS_TOKEN environment variable not set")
        print("Set it with: export PULUMI_ACCESS_TOKEN=<your-token>")
        sys.exit(1)

    return {
        "Authorization": f"token {token}",
        "Accept": "application/vnd.pulumi+8",
        "Content-Type": "application/json",
    }


def get_default_org() -> Optional[str]:
    """Get default organization from Pulumi CLI."""
    try:
        result = subprocess.run(
            ["pulumi", "org", "get-default"],
            capture_output=True,
            text=True,
            timeout=10,
        )
        if result.returncode == 0 and result.stdout.strip():
            org = result.stdout.strip()
            # Check if it looks like a valid cloud org (not a local path)
            if not org.startswith("/") and not org.startswith("file://"):
                return org
    except (subprocess.TimeoutExpired, FileNotFoundError):
        pass
    return None


def create_task(
    org: str,
    message: str,
    stack_name: Optional[str] = None,
    stack_project: Optional[str] = None,
    repo_name: Optional[str] = None,
    repo_org: Optional[str] = None,
    repo_forge: str = "github",
) -> str:
    """Create a new Neo task and return the task ID."""
    url = f"{BASE_URL}/{org}/tasks"

    entities_to_add = []

    if stack_name and stack_project:
        entities_to_add.append({
            "type": "stack",
            "name": stack_name,
            "project": stack_project,
        })

    if repo_name and repo_org:
        entities_to_add.append({
            "type": "repository",
            "name": repo_name,
            "org": repo_org,
            "forge": repo_forge,
        })

    payload = {
        "message": {
            "type": "user_message",
            "content": message,
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "entity_diff": {
                "add": entities_to_add,
                "remove": [],
            },
        }
    }

    response = requests.post(url, headers=get_headers(), json=payload)

    if response.status_code == 201:
        data = response.json()
        return data["taskId"]
    else:
        print(f"Error creating task: {response.status_code}")
        print(response.text)
        sys.exit(1)


def get_task(org: str, task_id: str) -> dict:
    """Get task metadata."""
    url = f"{BASE_URL}/{org}/tasks/{task_id}"
    response = requests.get(url, headers=get_headers())

    if response.status_code == 200:
        return response.json()
    else:
        print(f"Error getting task: {response.status_code}")
        print(response.text)
        sys.exit(1)


def list_tasks(org: str, page_size: int = 20) -> list:
    """List all tasks for an organization."""
    url = f"{BASE_URL}/{org}/tasks"
    params = {"pageSize": page_size}
    response = requests.get(url, headers=get_headers(), params=params)

    if response.status_code == 200:
        return response.json().get("tasks", [])
    else:
        print(f"Error listing tasks: {response.status_code}")
        print(response.text)
        sys.exit(1)


def get_events(
    org: str,
    task_id: str,
    continuation_token: Optional[str] = None,
) -> tuple[list, Optional[str]]:
    """Get task events with pagination."""
    url = f"{BASE_URL}/{org}/tasks/{task_id}/events"
    params = {"pageSize": 100}
    if continuation_token:
        params["continuationToken"] = continuation_token

    response = requests.get(url, headers=get_headers(), params=params)

    if response.status_code == 200:
        data = response.json()
        return data.get("events", []), data.get("continuationToken")
    else:
        print(f"Error getting events: {response.status_code}")
        print(response.text)
        return [], None


def send_approval(org: str, task_id: str, approval_request_id: str) -> bool:
    """Send approval confirmation for a pending request."""
    url = f"{BASE_URL}/{org}/tasks/{task_id}"

    payload = {
        "event": {
            "type": "user_confirmation",
            "approval_request_id": approval_request_id,
            "timestamp": datetime.now(timezone.utc).isoformat(),
        }
    }

    response = requests.post(url, headers=get_headers(), json=payload)
    return response.status_code == 202


def send_cancel(org: str, task_id: str) -> bool:
    """Send cancellation for a pending request."""
    url = f"{BASE_URL}/{org}/tasks/{task_id}"

    payload = {
        "event": {
            "type": "user_cancel",
            "timestamp": datetime.now(timezone.utc).isoformat(),
        }
    }

    response = requests.post(url, headers=get_headers(), json=payload)
    return response.status_code == 202


def send_message(org: str, task_id: str, message: str) -> bool:
    """Send a follow-up message to an existing task."""
    url = f"{BASE_URL}/{org}/tasks/{task_id}"

    payload = {
        "event": {
            "type": "user_message",
            "content": message,
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "entity_diff": {"add": [], "remove": []},
        }
    }

    response = requests.post(url, headers=get_headers(), json=payload)
    return response.status_code == 202


def format_event(event: dict) -> str:
    """Format an event for display."""
    event_type = event.get("type", "unknown")
    body = event.get("eventBody", {})
    inner_type = body.get("type", "")

    if event_type == "agentResponse":
        content = body.get("content", "")
        result = f"\n[Neo] {content}" if content else ""

        # Check for tool calls (including approval requests)
        tool_calls = body.get("tool_calls", [])
        for call in tool_calls:
            func = call.get("function", {})
            func_name = func.get("name", "")
            if func_name == "approval_request":
                try:
                    args = json.loads(func.get("arguments", "{}"))
                    req_id = args.get("approval_request_id", "")
                    desc = args.get("description", "")
                    result += f"\n[Approval Required] {desc}\n  Request ID: {req_id}"
                except json.JSONDecodeError:
                    pass
        return result if result else f"\n[agentResponse:{inner_type}]"
    elif event_type == "userInput":
        content = body.get("content", "")
        return f"\n[You] {content}"
    else:
        return f"\n[{event_type}:{inner_type}] {json.dumps(body, indent=2)}"


def find_pending_approval(events: list) -> Optional[str]:
    """Find pending approval request ID from events."""
    for event in reversed(events):
        if event.get("type") == "agentResponse":
            body = event.get("eventBody", {})
            tool_calls = body.get("tool_calls", [])
            for call in tool_calls:
                func = call.get("function", {})
                if func.get("name") == "approval_request":
                    try:
                        args = json.loads(func.get("arguments", "{}"))
                        return args.get("approval_request_id")
                    except json.JSONDecodeError:
                        pass
    return None


def poll_task(
    org: str,
    task_id: str,
    poll_interval: int = 5,
    max_wait: int = 600,
) -> None:
    """Poll a task for updates until completion or timeout."""
    print(f"Polling task {task_id}...")
    print(f"View in console: https://app.pulumi.com/{org}/neo/{task_id}")
    print("-" * 60)

    seen_events = set()
    all_events = []
    continuation_token = None
    start_time = time.time()
    last_status = None

    while True:
        # Check timeout
        elapsed = time.time() - start_time
        if elapsed > max_wait:
            print(f"\nTimeout after {max_wait} seconds. Task may still be running.")
            print(f"Continue with: python neo_task.py --org {org} --task-id {task_id}")
            break

        # Get task status
        task = get_task(org, task_id)
        status = task.get("status", "unknown")

        # Get new events
        events, new_token = get_events(org, task_id, continuation_token)

        for event in events:
            event_id = event.get("id")
            if event_id and event_id not in seen_events:
                seen_events.add(event_id)
                all_events.append(event)
                print(format_event(event))

        if new_token:
            continuation_token = new_token

        # Check for pending approval (detected from events, not status)
        pending_approval_id = find_pending_approval(all_events)

        # If status is idle, check if we're waiting for approval or truly done
        if status == "idle":
            if pending_approval_id:
                print(f"\n{'='*60}")
                print("Task is waiting for approval.")
                print(f"Approve: python neo_task.py --org {org} --task-id {task_id} --approve")
                print(f"Cancel:  python neo_task.py --org {org} --task-id {task_id} --cancel")
                break
            else:
                # Task finished processing
                print(f"\n{'='*60}")
                print("Task completed (idle)")
                break

        # If status changed from running to idle, task is done
        if last_status == "running" and status == "idle":
            print(f"\n{'='*60}")
            print("Task completed")
            break

        last_status = status
        time.sleep(poll_interval)


def main():
    # Validate token early
    token = os.environ.get("PULUMI_ACCESS_TOKEN")
    if not token:
        print("Error: PULUMI_ACCESS_TOKEN environment variable not set")
        print("Set it with: export PULUMI_ACCESS_TOKEN=<your-token>")
        sys.exit(1)

    parser = argparse.ArgumentParser(
        description="Pulumi Neo Task Manager",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Create a new task (with polling)
  %(prog)s --org myorg --message "Analyze my infrastructure"

  # Create task without polling (returns immediately - best for Claude Code)
  %(prog)s --org myorg --message "Analyze my infrastructure" --no-poll

  # Create task with JSON output (best for programmatic use)
  %(prog)s --org myorg --message "Analyze this" --no-poll --json

  # Create task with stack context
  %(prog)s --org myorg --message "Optimize this stack" \\
    --stack-name prod --stack-project my-infra --no-poll

  # List tasks
  %(prog)s --org myorg --list

  # Get current events for a task (single fetch, no polling)
  %(prog)s --org myorg --task-id TASK_ID --get-events

  # Poll an existing task for updates
  %(prog)s --org myorg --task-id TASK_ID

  # Approve a pending request
  %(prog)s --org myorg --task-id TASK_ID --approve

  # Cancel a pending request
  %(prog)s --org myorg --task-id TASK_ID --cancel
""",
    )

    parser.add_argument(
        "--org",
        help="Pulumi organization name (auto-detected if not specified)",
    )
    parser.add_argument(
        "--message", "-m",
        help="Message to send to Neo (creates new task if no task-id)",
    )
    parser.add_argument(
        "--task-id", "-t",
        help="Existing task ID to poll or respond to",
    )
    parser.add_argument(
        "--list", "-l",
        action="store_true",
        help="List existing tasks",
    )
    parser.add_argument(
        "--approve",
        action="store_true",
        help="Approve a pending request",
    )
    parser.add_argument(
        "--cancel",
        action="store_true",
        help="Cancel a pending request",
    )
    parser.add_argument(
        "--approval-id",
        help="Specific approval request ID (auto-detected from last event if not specified)",
    )

    # Entity context options
    parser.add_argument(
        "--stack-name",
        help="Stack name for context",
    )
    parser.add_argument(
        "--stack-project",
        help="Stack project for context (required with --stack-name)",
    )
    parser.add_argument(
        "--repo-name",
        help="Repository name for context",
    )
    parser.add_argument(
        "--repo-org",
        help="Repository organization for context (required with --repo-name)",
    )
    parser.add_argument(
        "--repo-forge",
        default="github",
        choices=["github", "gitlab", "bitbucket"],
        help="Repository forge (default: github)",
    )

    # Polling options
    parser.add_argument(
        "--poll-interval",
        type=int,
        default=5,
        help="Seconds between polls (default: 5)",
    )
    parser.add_argument(
        "--max-wait",
        type=int,
        default=600,
        help="Maximum seconds to wait (default: 600)",
    )
    parser.add_argument(
        "--no-poll",
        action="store_true",
        help="Create task and return immediately without polling (useful for CI/CD or Claude Code)",
    )
    parser.add_argument(
        "--get-events",
        action="store_true",
        help="Fetch and display current events once without polling",
    )
    parser.add_argument(
        "--json",
        action="store_true",
        help="Output results as JSON (useful for programmatic use)",
    )

    args = parser.parse_args()

    # Determine organization
    org = args.org
    if not org:
        org = get_default_org()
        if org:
            print(f"Using organization: {org}")
        else:
            print("Error: Could not detect organization.")
            print("Please specify with --org <organization>")
            print("")
            print("To set a default organization:")
            print("  pulumi org set-default <organization>")
            sys.exit(1)

    # Handle list command
    if args.list:
        tasks = list_tasks(org)
        if args.json:
            print(json.dumps({"tasks": tasks}, indent=2))
        elif not tasks:
            print("No tasks found.")
        else:
            print(f"{'ID':<40} {'Status':<15} {'Created'}")
            print("-" * 80)
            for task in tasks:
                print(f"{task.get('id', 'N/A'):<40} {task.get('status', 'N/A'):<15} {task.get('createdAt', 'N/A')}")
        return

    # Handle get-events command (single fetch, no polling)
    if args.get_events:
        if not args.task_id:
            print("Error: --task-id required for --get-events")
            sys.exit(1)
        task = get_task(org, args.task_id)
        events, _ = get_events(org, args.task_id)
        if args.json:
            print(json.dumps({"task": task, "events": events}, indent=2))
        else:
            print(f"Task: {args.task_id}")
            print(f"Status: {task.get('status', 'unknown')}")
            print(f"Console: https://app.pulumi.com/{org}/neo/{args.task_id}")
            print("-" * 60)
            for event in events:
                print(format_event(event))
        return

    # Handle approval/cancel
    if args.approve or args.cancel:
        if not args.task_id:
            print("Error: --task-id required for approval/cancel")
            sys.exit(1)

        approval_id = args.approval_id
        if not approval_id and args.approve:
            # Try to find approval ID from events
            events, _ = get_events(org, args.task_id)
            approval_id = find_pending_approval(events)

        if args.approve:
            if not approval_id:
                print("Error: Could not find approval request ID")
                print("Specify with --approval-id")
                sys.exit(1)
            if send_approval(org, args.task_id, approval_id):
                print("Approval sent. Continuing to poll...")
                poll_task(org, args.task_id, args.poll_interval, args.max_wait)
            else:
                print("Error sending approval")
                sys.exit(1)
        else:
            if send_cancel(org, args.task_id):
                print("Cancellation sent.")
            else:
                print("Error sending cancellation")
                sys.exit(1)
        return

    # Handle existing task polling
    if args.task_id and not args.message:
        poll_task(org, args.task_id, args.poll_interval, args.max_wait)
        return

    # Handle follow-up message to existing task
    if args.task_id and args.message:
        if send_message(org, args.task_id, args.message):
            print("Message sent. Polling for response...")
            poll_task(org, args.task_id, args.poll_interval, args.max_wait)
        else:
            print("Error sending message")
            sys.exit(1)
        return

    # Create new task
    if args.message:
        # Validate entity context combinations
        if args.stack_name and not args.stack_project:
            print("Error: --stack-project required with --stack-name")
            sys.exit(1)
        if args.repo_name and not args.repo_org:
            print("Error: --repo-org required with --repo-name")
            sys.exit(1)

        task_id = create_task(
            org=org,
            message=args.message,
            stack_name=args.stack_name,
            stack_project=args.stack_project,
            repo_name=args.repo_name,
            repo_org=args.repo_org,
            repo_forge=args.repo_forge,
        )

        console_url = f"https://app.pulumi.com/{org}/neo/{task_id}"

        if args.json:
            print(json.dumps({
                "taskId": task_id,
                "consoleUrl": console_url,
                "org": org,
            }, indent=2))
        else:
            print(f"Created task: {task_id}")
            print(f"Console: {console_url}")

        if not args.no_poll:
            poll_task(org, task_id, args.poll_interval, args.max_wait)
        return

    # No action specified
    parser.print_help()


if __name__ == "__main__":
    main()
