package commitlater

import "fmt"

type future struct {
	message string   // Message that describes this transaction.
	dir     string   // Basedir of the files
	files   []string // Names of the files
	display []string // Names as to be displayed to the user
}

// List of futures to be done in the future.
type List struct {
	items []*future
}

// Add queues up a future commit.
func (list *List) Add(message string, repobasedir string, files []string) {
	item := &future{
		message: message,
		dir:     repobasedir,
		files:   files,
	}
	list.items = append(list.items, item)
}

// Flush executes queued commits.
func (list *List) Flush(f func(string, string, []string) error) error {
	for _, fut := range list.items {
		err := f(fut.message, fut.dir, fut.files)
		if err != nil {
			return fmt.Errorf("commit files (%q) failed: %w", fut.files, err)
		}
	}
	return nil
}
