package commitlater

import (
	"fmt"
)

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

func sameDirs(l *List) bool {
	if len(l.items) <= 1 {
		return true
	}
	for _, k := range l.items[1:] {
		if k.dir != l.items[0].dir {
			return false
		}
	}
	return true
}

// Flush executes queued commits.
func (list *List) Flush(
	title string,
	fadd func([]string) error,
	fcommit func([]string, string, []string) error,
) error {

	// Just list the individual commit commands.
	if title == "" || len(list.items) < 2 || !sameDirs(list) {
		for _, fut := range list.items {
			err := fadd(fut.files)
			if err != nil {
				return fmt.Errorf("add files1 (%q) failed: %w", fut.files, err)
			}
			err = fcommit([]string{fut.message}, fut.dir, fut.files)
			if err != nil {
				return fmt.Errorf("commit files (%q) failed: %w", fut.files, err)
			}
		}
		return nil
	}

	// Create a long commit message.
	var m []string
	var f []string
	for _, fut := range list.items {
		err := fadd(fut.files)
		if err != nil {
			return fmt.Errorf("add files2 (%q) failed: %w", fut.files, err)
		}
		m = append(m, fut.message)
		f = append(f, fut.files...)
	}
	msg := []string{title}
	for _, mm := range m {
		msg = append(msg, "    * "+mm)
	}
	err := fcommit(msg, list.items[0].dir, f)
	if err != nil {
		return fmt.Errorf("commit files (%q) failed: %w", f, err)
	}

	return nil
}
