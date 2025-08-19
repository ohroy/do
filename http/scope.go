package dohttp

import (
	"github.com/samber/do/v2"
)

// ScopeTreeHTML generates an HTML page that displays the scope tree structure of the DI container.
// This function creates a visual representation of the scope hierarchy, showing the relationships
// between scopes and the services they contain.
//
// Parameters:
//   - basePath: The base URL path for the web interface
//   - injector: The injector instance to analyze
//   - scopeID: The ID of the specific scope to highlight (optional)
//
// Returns the HTML content as a string and any error that occurred during generation.
//
// The generated page includes:
//   - Visual representation of the scope hierarchy
//   - Service icons indicating their type and capabilities
//   - Navigation links between different views
//   - Interactive scope inspection links
//
// Example:
//
//	html, err := http.ScopeTreeHTML("/debug/di", injector, "")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Fprint(w, html)
func ScopeTreeHTML(basePath string, injector do.Injector, scopeID string) (string, error) {
	description := do.ExplainInjector(injector)

	return fromTemplate(
		`<!DOCTYPE html>
<html>
	<head>
		<title>Inspect scope tree - samber/do</title>
		<style>
		header {
			margin-bottom: 40px;
		}
		.scopes {
			margin-left: 10px;
			padding-top: 10px;
			padding-bottom: 10px;
			padding-left: 30px;
			border-left: 2px solid red;
		}
		.services {
			padding-top: 5px;
			padding-bottom: 5px;
		}
		</style>
	</head>
	<body>
		<h1>Scope description</h1>
		<small>
			Menu:
			<a href="{{.BasePath}}">Home</a>
			-
			<a href="{{.BasePath}}/scope">Scopes</a>
			-
			<a href="{{.BasePath}}/service">Services</a>
		</small>

		<header>
			<p>
				<b>Spec</b>:
				<br><br>
				😴 Lazy service
				<br>
				🔁 Eager service
				<br>
				🏭 Transient service
				<br>
				🔗 Service alias
				<br>
				🫀 Implements Healthchecker
				<br>
				🙅 Implements Shutdowner
			</p>
		</header>

		{{if .Scopes}}
			<ul class="scopes">
				{{range .Scopes}}
					<li class="scope">
						{{.}}
					</li>
				{{end}}
			</ul>
		{{end}}
	</body>
</html>`,
		map[string]any{
			"BasePath": basePath,
			"Scopes": mAp(description.DAG, func(item do.ExplainInjectorScopeOutput) string {
				return scopeTreeScopeToHTML(basePath, item)
			}),
		},
	)
}

func scopeTreeScopeToHTML(basePath string, description do.ExplainInjectorScopeOutput) string {
	html, _ := fromTemplate(
		`
			Scope:
			<a href="{{.BasePath}}/scope?scope_id={{.ScopeID}}">
				{{.ScopeName}}
			</a>

			{{if .Services}}
				<ul class="services">
					{{range .Services}}
						<li class="service">
							{{.}}
						</li>
					{{end}}
				</ul>
			{{end}}

			{{if .Scopes}}
				<ul class="scopes">
					{{range .Scopes}}
						<li class="scope">
							{{.}}
						</li>
					{{end}}
				</ul>
			{{end}}
		`,
		map[string]any{
			"BasePath":  basePath,
			"ScopeID":   description.ScopeID,
			"ScopeName": description.ScopeName,
			"Services": mAp(description.Services, func(item do.ExplainInjectorServiceOutput) string {
				return scopeTreeServiceToHTML(basePath, description.ScopeID, item)
			}),
			"Scopes": mAp(description.Children, func(item do.ExplainInjectorScopeOutput) string {
				return scopeTreeScopeToHTML(basePath, item)
			}),
		},
	)
	return html
}

func scopeTreeServiceToHTML(basePath string, scopeID string, description do.ExplainInjectorServiceOutput) string {
	featuresIcons := ""

	if description.IsHealthchecker {
		featuresIcons += " 🫀"
	}

	if description.IsShutdowner {
		featuresIcons += " 🙅"
	}

	html, _ := fromTemplate(
		`
			{{.ServiceTypeIcon}}
			<a href="{{.BasePath}}/service?scope_id={{.ScopeID}}&service_name={{.ServiceName}}">
				{{.ServiceName}}
			</a>
			{{.FeaturesIcons}}
		`,
		map[string]any{
			"BasePath":        basePath,
			"ScopeID":         scopeID,
			"ServiceName":     description.ServiceName,
			"ServiceTypeIcon": description.ServiceTypeIcon,
			"FeaturesIcons":   featuresIcons,
		},
	)
	return html
}
