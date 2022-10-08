package main

import (
	"fmt"
  "os"
  "strings"

  tm "github.com/buger/goterm"
  "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
  pointer int
  choices []string
  inputMode bool

  textInput textinput.Model
}


func initModel() model {
  ti := textinput.New()
	ti.Placeholder = "New entry here..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
  
  return model {
    choices: retrieve(),//[]string{"Find billy", "Laugh hysterically"},
    pointer: 0,
    inputMode: false,
    textInput: ti,
  }
}

func remove(i int, s []string) []string {
  if len(s) != 0 {
  copy(s[i:], s[i+1:])
  s[len(s)-1] = "" // remove element
  s = s[:len(s)-1]
  }
  return s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd){
  var cmd tea.Cmd
  
  switch msg := msg.(type) {
    case tea.KeyMsg:

      switch msg.String() {

        case "down":
          if m.pointer < len(m.choices) -1 {
            m.pointer += 1
          }

        case "up":
          if m.pointer > 0 {
            m.pointer -= 1
          }

        case "enter":
          switch m.inputMode {
            case false:
              m.choices = remove(m.pointer, m.choices)
              if m.pointer != 0 {
                m.pointer -= 1
              }
            case true:
              m.choices = append(m.choices, m.textInput.Value())
              m.textInput.SetValue("")
              m.inputMode = false
          }

        case "n":
         switch m.inputMode {
           case false:
              m.inputMode = true
              return m, nil
          }
          fallthrough
          //return m, cmd

        case "q":
          if m.inputMode == false {
            save(m.choices)
            return m, tea.Quit
          }
          fallthrough

        case "s":
          if m.inputMode == false {
            save(m.choices)
            return m, nil
          }
          fallthrough

        default:
          if m.inputMode == true {
          m.textInput, cmd = m.textInput.Update(msg)
        }
      } 
  }

  /**
  if msg.(type).String() != "+" {
    
  }
**/
  return m, cmd
}

func (m model) View() string {
  s := ""

  switch m.inputMode {
    case false:
      start := "Press 'n' for new entry!\n"
      for i, choice := range m.choices {
        p := "  "
        if i == m.pointer {
          p = ">>"
        }
        
        s += fmt.Sprintf("%s   %s\n", p, choice)  //choice + "\n"
      }
      if s == "" {
        s = "Nothing to do!\n"
      }
      s = start + s
      break
    
    case true:
      s += "New entry:\n" + m.textInput.View()
 }
  return s
}


func (m model) Init() tea.Cmd {
  return nil
}

func save(entries []string) {
  s := ""
  for _, entry := range entries {
    if len(entry) > 1 {
      s += entry + "\n"
    }
  }
  d1 := []byte(s[:len(s)-1])
    err := os.WriteFile("entries", d1, 0644)
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }
}

func retrieve() []string {
  s := []string{}
  dat, err := os.ReadFile("entries")
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  for _, entry := range strings.Split(string(dat), "\n") {
    if entry != "\n" {
      s = append(s, entry)
    }
  }
  return s
}


func clear() {
  tm.Clear()
 tm.MoveCursor(1, 1)
}


func main() {
  clear()
  p := tea.NewProgram(initModel())
	if err := p.Start(); err != nil {
		fmt.Println(err)
    os.Exit(1)
	}
}