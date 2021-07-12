package ac

import "github.com/kuhsinyv/dart"

const (
	FailState = -1
	RootState = 1
)

// Automaton 基于 DoubleArrayTrie 的 AC 自动机
type Automaton struct {
	trie   *dart.DoubleArrayTrie
	fail   []int
	output map[int][][]rune
}

// Result 匹配结果集
type Result struct {
	Pos  int
	Word []rune
}

// getF 返回 fail 中索引为 index 的元素
func (ac *Automaton) getF(index int) int {
	return ac.fail[index]
}

// setF 设置 fail 中索引为 inState 位置的值为 outState
func (ac *Automaton) setF(inState, outState int) {
	ac.fail[inState] = outState
}

// g 返回 out state
func (ac *Automaton) g(inState int, input rune) int {
	if inState == FailState {
		return RootState
	}

	t := inState + int(input) + dart.RootNodeBase
	if t >= len(ac.trie.Base) {
		if inState == RootState {
			return RootState
		}

		return FailState
	}

	if inState == ac.trie.Check[t] {
		return ac.trie.Base[t]
	}

	if inState == RootState {
		return RootState
	}

	return FailState
}

// Build 构建 AC 自动机
func (ac *Automaton) Build(patterns []string) error {
	var err error

	var trie *dart.LinkedListTrie

	d := new(dart.Dart)

	ac.trie, trie, err = d.Build(patterns)
	if err != nil {
		return err
	}

	ac.output = make(map[int][][]rune)
	for k, v := range d.Output {
		ac.output[k] = append(ac.output[k], v)
	}

	queue := make([]*dart.LinkedListTrieNode, 0)

	ac.fail = make([]int, len(ac.trie.Base))
	for _, child := range trie.Root.Children {
		ac.fail[child.Base] = dart.RootNodeBase
	}

	queue = append(queue, trie.Root.Children...)

	for {
		if len(queue) == 0 {
			break
		}

		node := queue[0]
		for _, c := range node.Children {
			if c.Base == dart.EndNodeBase {
				continue
			}

			inState := ac.getF(node.Base)
		setState:
			outState := ac.g(inState, c.Code-dart.RootNodeBase)

			if outState == FailState {
				inState = ac.getF(inState)

				goto setState
			}

			if _, ok := ac.output[outState]; ok {
				copyOutState := make([][]rune, 0)
				copyOutState = append(copyOutState, ac.output[outState]...)
				ac.output[c.Base] = append(copyOutState, ac.output[c.Base]...)
			}

			ac.setF(c.Base, outState)
		}

		queue = append(queue, node.Children...)
		queue = queue[1:]
	}

	return nil
}

// MultiPatternSearch 多模式匹配
func (ac *Automaton) MultiPatternSearch(content []rune) []*Result {
	terms := make([]*Result, 0)
	state := RootState

	for pos, c := range content {
	start:
		if ac.g(state, c) == FailState {
			state = ac.getF(state)

			goto start
		} else {
			state = ac.g(state, c)
			if v, ok := ac.output[state]; ok {
				for _, word := range v {
					term := &Result{
						Pos:  pos - len(word) + 1,
						Word: word,
					}
					terms = append(terms, term)
				}
			}
		}
	}

	return terms
}
