package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AddCoAuthor = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Add co-author on a commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("original")

		shell.SetConfig("user.email", "Bill@example.com")
		shell.SetConfig("user.name", "Bill Smith")

		shell.EmptyCommit("one")

		shell.NewBranch("other")

		shell.SetConfig("user.email", "John@example.com")
		shell.SetConfig("user.name", "John Smith")

		shell.EmptyCommit("two")

		shell.Checkout("original")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("BS").Contains("one").IsSelected(),
			)

		t.Views().Branches().
			Focus().
			Lines(
				Contains("original").IsSelected(),
				Contains("other"),
			).
			NavigateToLine(Contains("other")).
			PressEnter()

		// ensuring we get these commit authors as suggestions
		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("JS").Contains("two").IsSelected(),
				Contains("BS").Contains("one"),
			)

		t.Views().Commits().
			Focus().
			Press(keys.Commits.ResetCommitAuthor).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Amend commit attribute")).
					Select(Contains("Add co-author")).
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Contains("Add co-author")).
					SuggestionLines(
						Contains("Bill Smith"),
						Contains("John Smith"),
					).
					ConfirmSuggestion(Contains("John Smith"))
			})
		t.Views().Main().ContainsLines(
			Contains("Co-authored-by"),
		)
	},
})
