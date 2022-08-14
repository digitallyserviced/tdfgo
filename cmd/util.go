/*
    ┌ ┌┐
├- ┌┤ ├  ┬┐ ┌┐
│  ││ │  ││ ││
└┘ └┘ ┘  ┴┤ └┘
         ─┘
*/

package cmd

var fontsExamples = `
<fg=242># list fonts with impact in name</>
<primary>$</> <green>tdfgo</> <yellow>fonts</> impact

<fg=242># list fonts with much more information</>
<primary>$</> <green>tdfgo</> <yellow>fonts</> -v

<fg=242># list fonts with a preview output using default string "Preview"</>
<primary>$</> <green>tdfgo</> <yellow>fonts</> -p
`
var characterStatusFormat = `
<yellow>characters statuses:</>
%s
`
var watchExamples = `<fg=242># console clock</>
<primary>$</> <green>tdfgo</> <yellow>watch</> <cyan>-i 1s 'date "+%H:%m:%S"'</>
<primary>$</> <green>tdfgo</> <yellow>clock</>

<fg=242># print ONLINE or OFFLINE based off ping success </>
<primary>$</> <green>tdfgo</> <yellow>watch</> <cyan>'ping -c1 -w 1 google.com>/dev/null 2>&1 && echo -n ONLINE || echo -n OFFLINE'</>

<fg=242># print 1min load avg every 0.5 seconds </>
<primary>$</> <green>tdfgo</> <yellow>watch</> <cyan>'cat /proc/loadavg | cut -d" " -f 1'</>`

var usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

var helpTemplate = `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
