package main

import (
	"fmt"
	"strings"

	"github.com/dirien/yet-another-agent-harness/pkg/catalog"
	"github.com/spf13/cobra"
)

func skillsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skills",
		Short: "Discover and inspect available skills",
	}
	cmd.AddCommand(
		skillsListCmd(),
		skillsSearchCmd(),
		skillsInfoCmd(),
		skillsBundlesCmd(),
		skillsValidateCmd(),
	)
	return cmd
}

func loadCatalog(registryPath string) (*catalog.Catalog, error) {
	cat := catalog.DefaultCatalog()
	if registryPath != "" {
		extra, err := catalog.LoadRegistry(registryPath)
		if err != nil {
			return nil, err
		}
		cat = catalog.MergeCatalogs(cat, extra)
	}
	return cat, nil
}

func skillsListCmd() *cobra.Command {
	var (
		categoryFlag string
		tagFlag      string
		bundleFlag   string
		registryFlag string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all available skills from the catalog",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cat, err := loadCatalog(registryFlag)
			if err != nil {
				return err
			}
			var entries []catalog.CatalogEntry

			switch {
			case bundleFlag != "":
				entries = cat.ResolveBundle(bundleFlag)
				if entries == nil {
					return fmt.Errorf("unknown bundle: %s", bundleFlag)
				}
			case categoryFlag != "":
				entries = cat.ByCategory(catalog.Category(categoryFlag))
				if len(entries) == 0 {
					return fmt.Errorf("no skills in category: %s", categoryFlag)
				}
			case tagFlag != "":
				entries = cat.ByTag(tagFlag)
				if len(entries) == 0 {
					return fmt.Errorf("no skills with tag: %s", tagFlag)
				}
			default:
				entries = cat.Skills
			}

			printSkillTable(cmd, entries)
			return nil
		},
	}

	cmd.Flags().StringVar(&categoryFlag, "category", "", "Filter by category (language, infrastructure, security, devops, workflow, migration, architecture, testing)")
	cmd.Flags().StringVar(&tagFlag, "tag", "", "Filter by tag")
	cmd.Flags().StringVar(&bundleFlag, "bundle", "", "Show skills in a specific bundle")
	cmd.Flags().StringVar(&registryFlag, "registry", "", "Path or URL to an external skills-registry.json")
	return cmd
}

func skillsSearchCmd() *cobra.Command {
	var registryFlag string

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search skills by name, description, tags, or aliases",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cat, err := loadCatalog(registryFlag)
			if err != nil {
				return err
			}
			results := cat.Search(args[0])
			if len(results) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No skills matching %q\n", args[0])
				return nil
			}
			printSkillTable(cmd, results)
			return nil
		},
	}

	cmd.Flags().StringVar(&registryFlag, "registry", "", "Path or URL to an external skills-registry.json")
	return cmd
}

func skillsInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info <name>",
		Short: "Show detailed information about a skill",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cat := catalog.DefaultCatalog()
			entry := cat.Resolve(args[0])
			if entry == nil {
				return fmt.Errorf("unknown skill: %s", args[0])
			}
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Name:        %s\n", entry.Name)
			fmt.Fprintf(out, "Category:    %s\n", entry.Category)
			fmt.Fprintf(out, "Tags:        %s\n", strings.Join(entry.Tags, ", "))
			fmt.Fprintf(out, "Risk:        %s\n", entry.Risk)
			fmt.Fprintf(out, "Tier:        %s\n", entry.Tier)
			if entry.Uses != "builtin" {
				fmt.Fprintf(out, "Source:      %s\n", entry.Uses)
				if entry.Subpath != "" {
					fmt.Fprintf(out, "Subpath:     %s\n", entry.Subpath)
				}
			} else {
				fmt.Fprintf(out, "Source:      built-in\n")
			}
			if len(entry.Aliases) > 0 {
				fmt.Fprintf(out, "Aliases:     %s\n", strings.Join(entry.Aliases, ", "))
			}
			fmt.Fprintf(out, "Repo:        %s\n", entry.Repo)
			fmt.Fprintf(out, "Description: %s\n", entry.Description)

			bundles := cat.BundlesForSkill(entry.ID)
			if len(bundles) > 0 {
				names := make([]string, len(bundles))
				for i, b := range bundles {
					names[i] = b.ID
				}
				fmt.Fprintf(out, "Bundles:     %s\n", strings.Join(names, ", "))
			}
			return nil
		},
	}
}

func skillsBundlesCmd() *cobra.Command {
	var detail bool

	cmd := &cobra.Command{
		Use:   "bundles",
		Short: "List all skill bundles",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cat := catalog.DefaultCatalog()
			out := cmd.OutOrStdout()

			if detail {
				for _, b := range cat.Bundles {
					fmt.Fprintf(out, "%s (%d skills) — %s\n", b.ID, len(b.SkillIDs), b.Description)
					fmt.Fprintf(out, "  %s\n\n", strings.Join(b.SkillIDs, ", "))
				}
				return nil
			}

			fmt.Fprintf(out, "%-20s %6s  %s\n", "BUNDLE", "SKILLS", "DESCRIPTION")
			for _, b := range cat.Bundles {
				fmt.Fprintf(out, "%-20s %6d  %s\n", b.ID, len(b.SkillIDs), b.Description)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&detail, "detail", false, "Show member skills for each bundle")
	return cmd
}

func skillsValidateCmd() *cobra.Command {
	var nameFlag string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate that skill definitions can be fetched and are well-formed",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cat := catalog.DefaultCatalog()
			out := cmd.OutOrStdout()

			var entries []catalog.CatalogEntry
			if nameFlag != "" {
				entry := cat.Resolve(nameFlag)
				if entry == nil {
					return fmt.Errorf("unknown skill: %s", nameFlag)
				}
				entries = []catalog.CatalogEntry{*entry}
			} else {
				entries = cat.Skills
			}

			// Build a temporary catalog for validation.
			tmpCat := &catalog.Catalog{Skills: entries}
			results := catalog.ValidateAll(tmpCat, 8)

			passed, failed := 0, 0
			for _, r := range results {
				if r.OK {
					passed++
					if r.Size > 0 {
						fmt.Fprintf(out, "  ✓ %-28s content loaded (%.1fKB)\n", r.ID, float64(r.Size)/1024)
					} else {
						fmt.Fprintf(out, "  ✓ %-28s built-in\n", r.ID)
					}
				} else {
					failed++
					fmt.Fprintf(out, "  ✗ %-28s %s\n", r.ID, r.Error)
				}
			}

			fmt.Fprintf(out, "\n%d passed, %d failed\n", passed, failed)
			if failed > 0 {
				return fmt.Errorf("%d skill(s) failed validation", failed)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&nameFlag, "name", "", "Validate a single skill by name or alias")
	return cmd
}

func printSkillTable(cmd *cobra.Command, entries []catalog.CatalogEntry) {
	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "%-28s %-16s %-36s %s\n", "NAME", "CATEGORY", "TAGS", "TIER")
	for _, e := range entries {
		tags := strings.Join(e.Tags, ", ")
		if len(tags) > 34 {
			tags = tags[:31] + "..."
		}
		fmt.Fprintf(out, "%-28s %-16s %-36s %s\n", e.ID, e.Category, tags, e.Tier)
	}
	fmt.Fprintf(out, "\n%d skill(s)\n", len(entries))
}
