package clip

import (
	"fmt"
	"strings"
	"testing"
	// "github.com/duglin/clip"
)

func doTest(text string, width int, indent int, margin int) {
	fmt.Printf("w:%d i:%d m:%d  %.70q\n", width, indent, margin, text)
	fmt.Printf("12345678901234567890123456789012345678901234567890123456789012345678901234567890\n")
	fmt.Printf("%s\n", WrapWithArgs(text, width, indent, margin))
}

func trimLeading(text string) string {
	if strings.HasPrefix(text, "\n123") {
		text = text[4:]
		for text[0] != '\n' {
			text = text[1:]
		}
		text = text[1:]
	}
	return text
}

func TestWrap(t *testing.T) {
	type Test struct {
		Options  *WrapOptions
		Text     string
		Expected string
	}

	// Options perssist between tests to reduce repeating them
	tests := []Test{
		Test{},
		Test{
			Text:     "1",
			Expected: "1",
		},
		//             12345678901234567890
		Test{
			Text:     "hellohellohellohell",
			Expected: "hellohellohellohell",
		},
		Test{
			Text:     "hellohellohellohell\n",
			Expected: "hellohellohellohell",
		},
		Test{
			Text:     "hellohellohellohello",
			Expected: "hellohellohellohello",
		},
		Test{
			Text:     "hellohellohellohell ",
			Expected: "hellohellohellohell",
		},
		//             12345678901234567890
		Test{
			Text:     "hellohellohellohello ",
			Expected: "hellohellohellohello",
		},
		Test{
			Text:     "hellohellohellohello   ",
			Expected: "hellohellohellohello",
		},
		Test{
			Text:     "hellohellohellohello\nhello",
			Expected: "hellohellohellohello\nhello",
		},
		//             12345678901234567890
		Test{
			Text:     "hellohellohellohello hello",
			Expected: "hellohellohellohello\nhello",
		},
		Test{
			Text:     "hellohellohellohello     hello   ",
			Expected: "hellohellohellohello\nhello",
		},
		//             12345678901234567890
		Test{
			Text:     "hellohellohellohellohello",
			Expected: "hellohellohellohello\nhello",
		},

		Test{
			Text:     "hellohellohellohello\n",
			Expected: "hellohellohellohello",
		},
		Test{
			Text:     "hellohellohellohello\n\n",
			Expected: "hellohellohellohello\n",
		},
		//             12345678901234567890
		Test{
			Text:     "hellohellohellohello\nhello",
			Expected: "hellohellohellohello\nhello",
		},
		Test{
			Text:     "hellohellohellohello\n\nhello",
			Expected: "hellohellohellohello\n\nhello",
		},
		//             12345678901234567890
		Test{
			Text:     "hellohellohello hell",
			Expected: "hellohellohello hell",
		},
		Test{
			Text:     "hellohellohello hell\n",
			Expected: "hellohellohello hell",
		},
		Test{
			Text:     "hellohellohello hell\n\n",
			Expected: "hellohellohello hell\n",
		},
		Test{
			Text:     "hellohellohello hell\n ",
			Expected: "hellohellohello hell\n",
		},
		Test{
			Text:     "hellohellohello hell\n\n ",
			Expected: "hellohellohello hell\n\n",
		},
		Test{
			Text:     "hellohellohello hello",
			Expected: "hellohellohello\nhello",
		},
		//             12345678901234567890
		Test{
			Text:     "hello hello hello hello",
			Expected: "hello hello hello\nhello",
		},
		Test{
			Text:     "hello hello hello\nbye bye",
			Expected: "hello hello hello\nbye bye",
		},
		Test{
			Text:     "hello hello hello\n  bye bye",
			Expected: "hello hello hello\n  bye bye",
		},
		Test{
			Text:     "hello hello   hello\n  bye bye",
			Expected: "hello hello   hello\n  bye bye",
		},
		Test{
			Text:     "hello hello    hello\n  bye bye",
			Expected: "hello hello    hello\n  bye bye",
		},
		//             12345678901234567890
		Test{
			Text:     "\nhello hello hello\n  bye bye",
			Expected: "\nhello hello hello\n  bye bye",
		},

		// -----------------------------------------------------
		Test{
			Options: &WrapOptions{
				Width:  30,
				Indent: 5,
				Margin: 20,
			},
			Text: "-t,--timestamp\tA cool description of something",
			Expected: `
123456789012345678901234567890
     -t,--timestamp A cool
                    descriptio
                    n of
                    something`,
		},

		// -----------------------------------------------------
		Test{
			Text: "-t, --timestamp\tA cool description of something",
			Expected: `
123456789012345678901234567890
     -t, --timestamp
                    A cool
                    descriptio
                    n of
                    something`,
		},

		// -----------------------------------------------------
		Test{
			Text: "-t,  --timestamp\tA cool description of something",
			Expected: `
123456789012345678901234567890
     -t,  --timestamp
                    A cool
                    descriptio
                    n of
                    something`,
		},

		// -----------------------------------------------------
		Test{
			Text: "-t,  --timestamp\t\tA cool description of something",
			Expected: `
123456789012345678901234567890
     -t,  --timestamp

                    A cool
                    descriptio
                    n of
                    something`,
		},

		// -----------------------------------------------------
		Test{
			Options: &WrapOptions{
				Width:  20,
				Indent: 5,
				Margin: 10,
			},
			Text: "-t,  --timestamp\tA cool description of something",
			Expected: `
12345678901234567890
     -t,
          --timestam
          p
          A cool
          descriptio
          n of
          something`,
		},

		// -----------------------------------------------------
		Test{
			Text: "-t,\r--timestamp\tA cool description of something",
			Expected: `
12345678901234567890
     -t,
     --timestamp
          A cool
          descriptio
          n of
          something`,
		},

		// -----------------------------------------------------
		Test{
			Text: "-t,\r--timestamp\tA cool description\rof something",
			Expected: `
12345678901234567890
     -t,
     --timestamp
          A cool
          descriptio
          n
     of something`,
		},

		// -----------------------------------------------------
		Test{
			Text: "hello\t",
			Expected: `
12345678901234567890
     hello
`,
		},

		// -----------------------------------------------------
		Test{
			Text: "hi\t",
			Expected: `
12345678901234567890
     hi`,
		},

		// -----------------------------------------------------
		Test{
			Options: &WrapOptions{
				Width:  20,
				Indent: 5,
				Margin: 10,
			},
			Text: "hi\rhello\r",
			Expected: `
12345678901234567890
     hi
     hello
`,
		},

		// -----------------------------------------------------
		Test{
			Text: "-t,\r  --timestamp\tA cool description\r\rof something",
			Expected: `
12345678901234567890
     -t,
       --timestamp
          A cool
          descriptio
          n

     of something`,
		},

		// -----------------------------------------------------
		Test{
			Text: "-t,\r\t--time\tA cool description\r\rof something",
			Expected: `
12345678901234567890
     -t,
          --time
          A cool
          descriptio
          n

     of something`,
		},

		// -----------------------------------------------------
		Test{
			Text: "-t,\r\t--timestamp\tA cool description\r\rof something",
			Expected: `
12345678901234567890
     -t,
          --timestam
          p
          A cool
          descriptio
          n

     of something`,
		},

		// -----------------------------------------------------
		Test{
			Text: "-t,\r \t--timestamp\tA cool description\r\rof something",
			Expected: `
12345678901234567890
     -t,

          --timestam
          p
          A cool
          descriptio
          n

     of something`,
		},

		// -----------------------------------------------------
		Test{
			Options: &WrapOptions{
				Width:  20,
				Indent: 10,
				Margin: 5,
			},
			Text: "-t,\t\r--timestamp\tA cool description\r\rof something",
			Expected: `
12345678901234567890
          -t,
          --timestam
     p
     A cool
     description

          of
     something`,
		},

		// -----------------------------------------------------
		Test{
			Options: &WrapOptions{
				Width:  20,
				Indent: 5,
				Margin: 10,
			},
			Text: `
12345678901234567890

hellohell hello
  bye bye`,
			Expected: `
12345678901234567890

          hellohell
          hello
            bye bye`,
		},

		// -----------------------------------------------------
		Test{
			Text: `
12345678901234567890

hellohello hello
  bye bye`,
			Expected: `
12345678901234567890

          hellohello
          hello
            bye bye`,
		},

		// -----------------------------------------------------
		Test{
			Options: &WrapOptions{
				Width:             20,
				Indent:            5,
				Margin:            10,
				TrimNewlineSpaces: true,
			},
			Text: `
12345678901234567890
hellohello hello
  bye bye`,
			Expected: `
12345678901234567890
     hellohello
          hello
          bye bye`,
		},

		// -----------------------------------------------------
		Test{
			Options: &WrapOptions{
				Width:             14,
				Indent:            6,
				Margin:            4,
				TrimNewlineSpaces: false,
				Justify:           true,
			},
			Text: `
12345678901234567890
he lohe lo he lo hello hello`,
			Expected: `
12345678901234567890
      he  lohe
    lo  he  lo
    hello
    hello`,
		},

		// -----------------------------------------------------
		Test{
			Text: `
12345678901234567890
hel hi bye
   hell bye  hi bye
  hello`,
			Expected: `
12345678901234567890
      hel   hi
    bye
       hell
    bye     hi
    bye
      hello`,
		},
	}

	doTest("Hello world", 80, 0, 0)
	doTest("Hello world goodbye friends", 10, 0, 0)
	doTest("Hello world goodbye friends", 5, 0, 0)
	doTest("Hello world goodbye friends", 6, 0, 0)
	doTest("Hello\nworld goodbye friends", 5, 0, 0)
	doTest("Hello\nworld goodbye friends", 6, 0, 0)
	doTest("Hello\nworld goodbye friends", 1, 0, 0)
	doTest("Hello\n\nworld goodbye friends", 1, 0, 0)
	doTest("Hello\nworld goodbye friends", 2, 0, 0)
	doTest("Hello\n\nworld goodbye friends", 2, 0, 0)
	doTest("Hello world goodbye friends", 6, 0, 4)
	doTest("Hello world goodbye friends", 6, 3, 1)
	doTest("-t, --time\nA really cool description that's kind of long but htat should be ok", 20, 5, 10)
	doTest("-t, --time\tA really cool description that's kind of long but htat should be ok", 30, 5, 20)
	doTest("-t, --timestamp\tA really cool description that's kind of long but htat should be ok", 30, 5, 20)
	doTest("-t, --timestam\tA really cool description that's kind of long but that sould be ok", 80, 5, 20)
	fmt.Printf("Done\n")

	opts := NewWrapOptions()
	opts.Width = 20

	for _, test := range tests {
		if test.Options != nil {
			opts = test.Options
		}
		if opts.Indent < 0 {
			continue
		}
		result := opts.Wrap(trimLeading(test.Text))
		expected := trimLeading(test.Expected)
		line := strings.Repeat("1234567890", 1+opts.Width/10)[:opts.Width]
		if result != expected {
			fmt.Printf("=====================\n")
			fmt.Printf("Text: %q\n", trimLeading(test.Text))
			ShowDebug = true
			result := opts.Wrap(trimLeading(test.Text))
			t.Errorf("Failed:\n  opts: %#v\ntext:\n%q<<\n"+
				"%s\n"+
				"got:\n%s<<\nexpected:\n%s<<\n",
				opts, test.Text, line,
				strings.ReplaceAll(result, " ", "."),
				strings.ReplaceAll(expected, " ", "."))
		}
	}
}
