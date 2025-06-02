package internal

import (
    "encoding/json"
    "os"
    "os/user"
    "path/filepath"
)

type PersistentState struct {
    SelectedIndex int            `json:"selected_index"`
    Solutions     map[int]string `json:"solutions"`
}

func getStateFilePath() string {
    usr, _ := user.Current()
    configDir := filepath.Join(usr.HomeDir, ".ts-koans")
    os.MkdirAll(configDir, 0700)
    return filepath.Join(configDir, "state.json")
}

func SaveState(state PersistentState) error {
    path := getStateFilePath()
    data, _ := json.MarshalIndent(state, "", "  ")
    return os.WriteFile(path, data, 0600)
}

func LoadState() (PersistentState, error) {
    var state PersistentState
    path := getStateFilePath()
    data, err := os.ReadFile(path)
    if err != nil {
        state.Solutions = make(map[int]string)
        return state, nil
    }
    json.Unmarshal(data, &state)
    if state.Solutions == nil {
        state.Solutions = make(map[int]string)
    }
    return state, nil
}
