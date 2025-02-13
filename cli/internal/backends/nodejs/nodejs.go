package nodejs

import (
	"fmt"
	"io/ioutil"
	"turbo/internal/api"
	"turbo/internal/fs"
	"turbo/internal/util"

	"gopkg.in/yaml.v3"
)

// nodejsPatterns is the FilenamePatterns value for NodejsBackend.
var nodejsPatterns = []string{"*.js", ".mjs", "*.ts", "*.jsx", "*.tsx"}

var NodejsYarnBackend = api.LanguageBackend{
	Name:             "nodejs-yarn",
	Specfile:         "package.json",
	Lockfile:         "yarn.lock",
	FilenamePatterns: nodejsPatterns,
	GetWorkspaceGlobs: func() ([]string, error) {
		pkg, err := fs.ReadPackageJSON("package.json")
		if err != nil {
			return nil, fmt.Errorf("package.json: %w", err)
		}
		if len(pkg.Workspaces) == 0 {
			return nil, fmt.Errorf("package.json: no workspaces found. Turborepo requires Yarn workspaces to be defined in the root package.json")
		}
		return pkg.Workspaces, nil
	},
	GetPackageDir: func() string {
		return "node_modules"
	},
	GetRunCommand: func() []string {
		return []string{"yarn", "run"}
	},
	Detect: func(cwd string, pkg *fs.PackageJSON, backend *api.LanguageBackend) (bool, error) {
		if pkg.PackageManager != "" {
			packageManager, version := util.GetPackageManagerAndVersion(pkg.PackageManager)

			if packageManager != "yarn" {
				return false, nil
			}

			isBerry, err := util.IsBerry(cwd, version)
			if err != nil {
				return false, fmt.Errorf("could not determine yarn version (v1 or berry): %w", err)
			}

			if !isBerry {
				return true, nil
			}
		}

		return false, nil
	},
}

var NodejsBerryBackend = api.LanguageBackend{
	Name:             "nodejs-berry",
	Specfile:         "package.json",
	Lockfile:         "yarn.lock",
	FilenamePatterns: nodejsPatterns,
	GetWorkspaceGlobs: func() ([]string, error) {
		pkg, err := fs.ReadPackageJSON("package.json")
		if err != nil {
			return nil, fmt.Errorf("package.json: %w", err)
		}
		if len(pkg.Workspaces) == 0 {
			return nil, fmt.Errorf("package.json: no workspaces found. Turborepo requires Yarn workspaces to be defined in the root package.json")
		}
		return pkg.Workspaces, nil
	},
	GetPackageDir: func() string {
		return "node_modules"
	},
	GetRunCommand: func() []string {
		return []string{"yarn", "run"}
	},
	Detect: func(cwd string, pkg *fs.PackageJSON, backend *api.LanguageBackend) (bool, error) {
		if pkg.PackageManager != "" {
			packageManager, version := util.GetPackageManagerAndVersion(pkg.PackageManager)

			if packageManager != "yarn" {
				return false, nil
			}

			isBerry, err := util.IsBerry(cwd, version)
			if err != nil {
				return false, fmt.Errorf("could not determine yarn version (v1 or berry): %w", err)
			}

			if isBerry {
				isNMLinker, err := util.IsNMLinker(cwd)
				if err != nil {
					return false, fmt.Errorf("could not determine if yarn is using `nodeLinker: node-modules`: %w", err)
				} else if !isNMLinker {
					return false, fmt.Errorf("only yarn v2/v3 with `nodeLinker: node-modules` is supported at this time")
				}

				return true, nil
			}
		}

		return false, nil
	},
}

// PnpmWorkspaces is a representation of workspace package globs found
// in pnpm-workspace.yaml
type PnpmWorkspaces struct {
	Packages []string `yaml:"packages,omitempty"`
}

var NodejsPnpmBackend = api.LanguageBackend{
	Name:             "nodejs-pnpm",
	Specfile:         "package.json",
	Lockfile:         "pnpm-lock.yaml",
	FilenamePatterns: nodejsPatterns,
	GetWorkspaceGlobs: func() ([]string, error) {
		bytes, err := ioutil.ReadFile("pnpm-workspace.yaml")
		if err != nil {
			return nil, fmt.Errorf("pnpm-workspace.yaml: %w", err)
		}
		var pnpmWorkspaces PnpmWorkspaces
		if err := yaml.Unmarshal(bytes, &pnpmWorkspaces); err != nil {
			return nil, fmt.Errorf("pnpm-workspace.yaml: %w", err)
		}

		if len(pnpmWorkspaces.Packages) == 0 {
			return nil, fmt.Errorf("pnpm-workspace.yaml: no packages found. Turborepo requires PNPM workspaces and thus packages to be defined in the root pnpm-workspace.yaml")
		}

		return pnpmWorkspaces.Packages, nil
	},
	GetPackageDir: func() string {
		return "node_modules"
	},
	GetRunCommand: func() []string {
		return []string{"pnpm", "run"}
	},
	Detect: func(_ string, pkg *fs.PackageJSON, backend *api.LanguageBackend) (bool, error) {
		if pkg.PackageManager != "" {
			packageManager, _ := util.GetPackageManagerAndVersion(pkg.PackageManager)

			if packageManager == "pnpm" {
				return true, nil
			}
		}

		return false, nil
	},
}

var NodejsNpmBackend = api.LanguageBackend{
	Name:             "nodejs-npm",
	Specfile:         "package.json",
	Lockfile:         "package-lock.json",
	FilenamePatterns: nodejsPatterns,
	GetWorkspaceGlobs: func() ([]string, error) {
		pkg, err := fs.ReadPackageJSON("package.json")
		if err != nil {
			return nil, fmt.Errorf("package.json: %w", err)
		}
		if len(pkg.Workspaces) == 0 {
			return nil, fmt.Errorf("package.json: no workspaces found. Turborepo requires NPM workspaces to be defined in the root package.json")
		}
		return pkg.Workspaces, nil
	},
	GetPackageDir: func() string {
		return "node_modules"
	},
	GetRunCommand: func() []string {
		return []string{"npm", "run"}
	},
	Detect: func(_ string, pkg *fs.PackageJSON, backend *api.LanguageBackend) (bool, error) {
		if pkg.PackageManager != "" {
			packageManager, _ := util.GetPackageManagerAndVersion(pkg.PackageManager)

			if packageManager == "npm" {
				return true, nil
			}
		}

		return false, nil
	},
}
