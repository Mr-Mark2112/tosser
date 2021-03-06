package tosser

import "sync"

//processingCache - используется для отслеживания уже обрабатываемых файлов

type ProcessingCacher interface {
	add()
	del()
	check()
}
type processingCache struct {
	sync.RWMutex
	cache map[string]int
}

func (p *processingCache) add(fullSrcPath string) {
	p.Lock()
	p.cache[fullSrcPath] = 1
	p.Unlock()
}

func (p *processingCache) del(fullSrcPath string) {
	p.Lock()
	delete(p.cache, fullSrcPath)
	p.Unlock()
}
func (p *processingCache) check(fullSrcPath string) bool {
	p.Lock()
	_, ok := p.cache[fullSrcPath]
	p.Unlock()
	return ok
}

//NewProcessingCache - создает и возвращает ссылку на processingCache
func NewProcessingCache() *processingCache {
	p := new(processingCache)
	p.cache = make(map[string]int)
	return p
}
