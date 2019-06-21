package files

import (
	"fmt"
	"github.com/jlcheng/forget/db"
	"github.com/jlcheng/forget/orgmode"
	"github.com/jlcheng/forget/trace"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

func RebuildIndex(atlas *db.Atlas, dirs []string) error {
	walkf := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			trace.Warn(fmt.Sprintf("cannot index [%v]: %v", path, err))
			return err
		}

		if info.IsDir() {
			return nil
		}

		if db.FilterFile(path, info) {
			notes, err := ParseFile(path)
			if err != nil {
				trace.Warn(fmt.Sprintf("cannot parse [%v]: %v", path, err))
				return err
			}
			for _, note := range notes {
				err = atlas.Enqueue(note)
				if err != nil {
					trace.Warn(fmt.Sprintf("cannot index [%v]: %v", note.ID, err))
				}
			}
		}

		return nil
	}

	for _, dir := range dirs {
		err := filepath.Walk(dir, walkf)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func ParseFile(path string) ([]db.Note, error) {
	ret := make([]db.Note, 0)
	fi, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("cannot stat %v", path))
	}

	body, err := readFileAsString(path)
	if err != nil {
		return nil, err
	}

	if filepath.Ext(path) == ".org" {
		doc, err := orgmode.ParseOrgDoc(body)
		if err != nil {
			return nil, err
		}
		for _, n := range doc.Nodes {
			if n.Type == orgmode.NodeTypeHeading {
				continue
			}

			if n.Type == orgmode.NodeTypeText {
				ret = append(ret, db.Note{
					ID:         fmt.Sprintf("%v:%v", path, n.Heading),
					Body:       n.TextValue(),
					Title:      n.Heading,
					AccessTime: fi.ModTime().Unix(),
				})
			} else if n.Type == orgmode.NodeTypeDescItem {
				di := n.DescItemValue()
				ret = append(ret, db.Note{
					ID:         fmt.Sprintf("%v:%v:%v", path, n.Heading, di.Term),
					Body:       di.Desc,
					Title:      di.Term,
					AccessTime: fi.ModTime().Unix(),
				})
			} else {
				trace.Warn(fmt.Sprintf("unexpected node type %v in %v", n.Type, path))
			}
		}
	} else {
		ret = append(ret, db.Note{
			ID:         path,
			Body:       body,
			Title:      fi.Name(),
			AccessTime: fi.ModTime().Unix(),
		})
	}

	return ret, nil
}

func readFileAsString(fileName string) (string, error) {
	s, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("cannot open %v", fileName))
	}
	return string(s), nil
}
