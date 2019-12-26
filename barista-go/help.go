package barista

import (
	"fmt"
	"strings"

	"github.com/Necroforger/dgwidgets"

	"github.com/bwmarrin/discordgo"
)

type flag struct {
	name      string
	longForm  string
	shortForm string
	desc      string
}

type arg struct {
	name string
	desc string
}

func cmdEmbedWithArgs(cmd string, desc string, args []arg, flags []flag) *discordgo.MessageEmbed {
	embed := NewEmbed()
	title := cmd
	for _, val := range args {
		title = title + " [" + val.name + "] "
	}
	for _, val := range flags {
		title = title + fmt.Sprintf(" [ --%s/-%s %s ] ", val.longForm, val.shortForm, val.name)
	}
	embed.SetTitle(title)
	embed.SetDescription(desc)
	for _, val := range args {
		embed.AddField(val.name, val.desc, true)
	}
	for _, val := range flags {
		embed.AddField(fmt.Sprintf("--%s/-%s", val.longForm, val.shortForm), val.desc, true)
	}
	embed.SetColor(0xc12fb7)
	return embed.MessageEmbed
}

func cmdEmbed(cmd string, desc string) *discordgo.MessageEmbed {
	embed := NewEmbed()
	embed.SetTitle(cmd)
	embed.SetDescription(desc)
	embed.SetColor(0xc12fb7)
	return embed.MessageEmbed
}

const profilehelp string = `# Syntax: sudo profile --flag value
( --user | -u )
	Get the user specified.
( --set-desktop-environment | -w )
	Set your desktop environment or window manager.
( --set-distro | -d )
	Set your distro.
( --set-shell | -s )
	Set your command line shell.
( --set-editor | -e )
	Set your editor.
( --set-languages | -p )
	Set your programming languages.
( --set-blurb | -b )
	Set your profile blurb.
( --set-screenshot | -i )
	Set your screenshot to a PNG.`

const gsettingshelp string = `# Syntax: sudo gsettings [--action] [--schema schema] [--setting setting] [--value value]
( -g | --get )
	Get the value of SCHEMA > SETTING
( -s | --set )
	Set the value of SCHEMA > SETTING.
( -l | --list-schemas )
	List schemas.
( -k | --list-settings )
	List settings in SCHEMA.`

const repoqueryhelp string = `# Syntax: dnf repoquery [--flag value] --distro distro
( -f | --file )
	Get packages that have FILE
( --whatconflicts )
	Get packages that conflict PKG.
( --whatprovides )
	Get packages that provide CAPABILITY.
( --whatobsoletes)
	Get packages that obsolete PKG.
( --whatrecommends )
	Get packages that recommend PKG.
( --whatenhances )
	Get packages that enhance PKG.
( --whatsupplements )
	Get packages that supplement PKG.
( --whatsuggests )
	Get packages that supplement PKG.

( -l | --list )
	List files provided by PKG. Will override other flags.
( -n | --no-details )
	Only list package names without details.`

const msgtags string = `SR1234		openSUSE Build Service Submit Requests
FEDORA-*		Fedora Bodhi Updates`

func Help(s *discordgo.Session, cmd *LexedCommand) {
	page := dgwidgets.NewPaginator(s, cmd.CommandMessage.ChannelID)
	page.Add(
		cmdEmbed("sudo help", "This command"),
	)
	page.Add(
		cmdEmbedWithArgs(
			"dnf search", "Search through distro repositories",
			[]arg{arg{name: "query", desc: "The search to look for"}},
			[]flag{flag{longForm: "distro", shortForm: "d", desc: "The distro's repositories to search.", name: "distro"}},
		),
	)
	page.Add(
		cmdEmbed("dnf repoquery", "Query distro repos"+"```dsconfig\n"+repoqueryhelp+"\n```"),
	)
	page.Add(
		cmdEmbedWithArgs(
			"sudo ddg", "Get answers from DuckDuckGo. Note: due to issues with DDG, answers may not be relevant or even existent.",
			[]arg{arg{name: "query", desc: "The search to look for"}},
			[]flag{},
		),
	)
	page.Add(
		cmdEmbed("sudo profile", "Look at user profiles"+"```dsconfig\n"+profilehelp+"\n```"),
	)
	page.Add(
		cmdEmbed("sudo about", "About Barista"),
	)
	bztags := []string{}
	for _, bugzilla := range BugzillaInstances {
		bztags = append(bztags, fmt.Sprintf("%s\t\t%s", strings.Join(bugzilla.Matches, ", "), bugzilla.Name))
	}
	page.Add(
		cmdEmbed("Message Tags", "```"+msgtags+"\n"+strings.Join(bztags, "\n")+"```"),
	)
	cmd.PaginatorPageName = "Command"
	cmd.SendPaginator(page)
}
