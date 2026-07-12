package detection

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kdse/runtime/internal/types"
)

type Detector struct {
	repoPath string
}

func NewDetector(repoPath string) *Detector {
	return &Detector{repoPath: repoPath}
}

func (d *Detector) Detect() (*types.Repository, error) {
	name := filepath.Base(d.repoPath)
	isGit := d.isGitRepo()

	artifacts := d.detectArtifacts()

	phase := d.determinePhase(artifacts)

	return &types.Repository{
		Path:      d.repoPath,
		Name:      name,
		Phase:     string(phase),
		Artifacts: artifacts,
		IsGitRepo: isGit,
	}, nil
}

func (d *Detector) isGitRepo() bool {
	gitPath := filepath.Join(d.repoPath, ".git")
	info, err := os.Stat(gitPath)
	return err == nil && info.IsDir()
}

func (d *Detector) detectArtifacts() []string {
	var artifacts []string

	patterns := map[string][]string{
		"docs/":              {".md"},
		"src/":               {".go", ".py", ".js", ".ts", ".java", ".rs", ".cpp", ".c"},
		"tests/":             {".test.", ".spec.", "_test.", "-test."},
		"requirements/":       {".txt", ".lock", ".toml", ".yaml", ".yml"},
		"architecture/":      {".md", ".drawio", ".puml"},
		"commands/":          {".sh", ".bash"},
	}

	artifactTypes := []string{
		"README.md",
		"docs/",
		"src/",
		"tests/",
		"requirements/",
		"architecture/",
		"commands/",
		"SPEC.md",
		"ARCHITECTURE.md",
		"CONTRIBUTING.md",
		"LICENSE",
		"package.json",
		"go.mod",
		"requirements.txt",
		" Cargo.toml",
		"pom.xml",
	}

	seen := make(map[string]bool)

	filepath.Walk(d.repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" || name == "__pycache__" {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, err := filepath.Rel(d.repoPath, path)
		if err != nil {
			return nil
		}

		if seen[relPath] {
			return nil
		}

		for _, pattern := range artifactTypes {
			if strings.HasSuffix(pattern, "/") {
				if strings.HasPrefix(relPath, pattern) {
					dir := filepath.Dir(relPath)
					if !seen[dir] {
						artifacts = append(artifacts, dir)
						seen[dir] = true
					}
					return nil
				}
			} else if relPath == pattern {
				artifacts = append(artifacts, relPath)
				seen[relPath] = true
				return nil
			}
		}

		ext := filepath.Ext(info.Name())
		for category, extensions := range patterns {
			for _, e := range extensions {
				if ext == e {
					dir := filepath.Dir(relPath)
					if strings.HasPrefix(dir, category) || dir == strings.TrimSuffix(category, "/") {
						if !seen[category] {
							artifacts = append(artifacts, category)
							seen[category] = true
						}
					}
				}
			}
		}

		return nil
	})

	if len(artifacts) == 0 {
		artifacts = append(artifacts, "empty-repository")
	}

	return artifacts
}

func (d *Detector) determinePhase(artifacts []string) types.EngineeringPhase {
	has := func(items []string, target string) bool {
		for _, item := range items {
			if item == target || strings.HasPrefix(item, target) {
				return true
			}
		}
		return false
	}

	score := 0

	if has(artifacts, "README.md") || has(artifacts, "SPEC.md") {
		score += 2
	}
	if has(artifacts, "docs/") {
		score += 2
	}
	if has(artifacts, "architecture/") || has(artifacts, "ARCHITECTURE.md") {
		score += 3
	}
	if has(artifacts, "src/") {
		score += 2
	}
	if has(artifacts, "tests/") {
		score += 2
	}
	if has(artifacts, "requirements/") || has(artifacts, "commands/") {
		score += 2
	}
	if d.isGitRepo() {
		score += 1
	}

	switch {
	case score >= 10:
		return types.PhaseValidated
	case score >= 8:
		return types.PhaseUsable
	case score >= 5:
		return types.PhaseStructured
	case score >= 3:
		return types.PhaseDefined
	default:
		return types.PhaseConcept
	}
}
