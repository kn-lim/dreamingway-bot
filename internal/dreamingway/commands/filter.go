package commands

import "fmt"

// FilterCommands selects the commands that have a name in the names list.
// If the names list is empty, FilterCommands returns an empty map. The caller then syncs no commands.
// If a name has no command in the map, FilterCommands returns an error.
func FilterCommands(commands map[string]Command, names []string) (map[string]Command, error) {
	filtered := make(map[string]Command, len(names))

	for _, name := range names {
		found := false
		for key, cmd := range commands {
			if cmd.Command.CommandName() == name {
				filtered[key] = cmd
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("command not available in this scope: %q", name)
		}
	}

	return filtered, nil
}
