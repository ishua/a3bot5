package domain_test

import (
	"testing"

	"github.com/ishua/a3bot5/fsnotes/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestGitFile_TextClean(t *testing.T) {
	tests := []struct {
		name         string
		texts        []string
		expectedText string
	}{
		{
			name:         "easy test",
			texts:        []string{"note", "diary", "add", "5bx", "4", "5", "6", "7"},
			expectedText: "4 5 6 7",
		},
		{
			name:         "easy test1",
			texts:        []string{"note", "diary", "add", "5bx", "2"},
			expectedText: "2",
		},
		{
			name:         "easy test2",
			texts:        []string{"note", "diary", "add", "5bx", "4", "5\n6", "7"},
			expectedText: "4 5 6 7",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultsText := domain.TextClean(tt.texts)
			assert.Equal(t, tt.expectedText, resultsText)
		})
	}

}

func Test_ParceCommand(t *testing.T) {
	tests := []struct {
		name            string
		msg             domain.Msg
		expectedNote    domain.Note
		expectedErrText string
	}{
		{
			name:            "add 5bx",
			msg:             domain.Msg{Text: "/note diary add 5bx 2"},
			expectedNote:    domain.Note{Theme: "5bx", Text: "2"},
			expectedErrText: "",
		},
		{
			name:            "add entry",
			msg:             domain.Msg{Text: "/note diary add entry some entry with some word"},
			expectedNote:    domain.Note{Theme: "entry", Text: "some entry with some word"},
			expectedErrText: "",
		},
		{
			name:            "add entry with synonyms",
			msg:             domain.Msg{Text: "note d a e some entry with some word"},
			expectedNote:    domain.Note{Theme: "entry", Text: "some entry with some word"},
			expectedErrText: "",
		},
		{
			name:            "Wrong fsnotes command",
			msg:             domain.Msg{Text: "/note diary"},
			expectedNote:    domain.Note{},
			expectedErrText: "Wrong fsnotes command",
		},
		{
			name:            "Wrong fsnotes command with 5 word",
			msg:             domain.Msg{Text: "/note diary add"},
			expectedNote:    domain.Note{},
			expectedErrText: "Wrong fsnotes command",
		},
		{
			name:            "Not working",
			msg:             domain.Msg{Text: "/note diary list"},
			expectedNote:    domain.Note{},
			expectedErrText: "Not working",
		},
		{
			name:            "fsnotes doesn't have this command",
			msg:             domain.Msg{Text: "/note diary some"},
			expectedNote:    domain.Note{},
			expectedErrText: "fsnotes doesn't have this command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			note, err := tt.msg.ParseCommand()
			var errText string
			if err != nil {
				errText = err.Error()
			}
			assert.Equal(t, tt.expectedNote, note)
			assert.Equal(t, tt.expectedErrText, errText)
		})
	}
}
