package main

import (
  "charm.land/lipgloss/v2"
)

var (
  FocusStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
  InputStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
  ResultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
  ErrorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
  HelpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))

  NestedColors = []lipgloss.Style{
    lipgloss.NewStyle().Foreground(lipgloss.Color("205")), // pink
    lipgloss.NewStyle().Foreground(lipgloss.Color("75")),  // cyan
    lipgloss.NewStyle().Foreground(lipgloss.Color("197")), // magenta
    lipgloss.NewStyle().Foreground(lipgloss.Color("33")),  // blue
    lipgloss.NewStyle().Foreground(lipgloss.Color("130")), // orange
  }

  BaseTableStyle = lipgloss.NewStyle().
      BorderStyle(lipgloss.NormalBorder()).
      BorderForeground(lipgloss.Color("240"))
)
