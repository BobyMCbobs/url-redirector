// This program is free software: you can redistribute it and/or modify
// it under the terms of the Affero GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the Affero GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package types ...
// project types
package types

// Routes ...
// path to destination mapping
type Routes map[string]string

// RouteHost ...
// additional paths for mapping
type RouteHost struct {
	Routes   Routes `yaml:"routes"`
	Root     string `yaml:"root"`
	Wildcard string `yaml:"wildcard"`
}

// RouteHosts ...
// group routes by host
type RouteHosts map[string]RouteHost
