package main

import (
	"errors"
	"fmt"
	"math"
	"syscall"
	"unsafe"
)

const (
	MemoryHeap = 32768
	PageSize   = 4096

	FreePage    = 0
	DividedPage = 1
	MultiPage   = 2
)

type PageAllocator struct {
	heapSize     uint
	heap         []byte
	pages        []*PageHeader
}

type PageHeader struct {
	state        uint
	size         uint
	adr          uint
	pointer      unsafe.Pointer
	blocksAmount uint
	blocks       []BlockHeader
}

type BlockHeader struct {
	adr  uint
	next unsafe.Pointer
}

func (p *PageAllocator) Init(sizeHeap uint) (err error) {
	if sizeHeap < PageSize {
		return errors.New("size of heap is too small")
	}
	p.heapSize = PageSize * (sizeHeap / PageSize)

	p.heap, err = syscall.Mmap(-1, 0, int(p.heapSize), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS)
	if err != nil {
		return err
	}
	p.startPointer = unsafe.Pointer(&p.heap[0])
	for i := 0; i < (int(p.heapSize) / PageSize); i++ {
		p.pages = append(p.pages, &PageHeader{
			state:   FreePage,
			pointer: unsafe.Pointer(&p.heap[PageSize*i]),
			adr:     uint(PageSize * i),
		})
		p.freePages = append(p.freePages, p.pages[i])

	}
	return nil
}

func (p *PageAllocator) mem_alloc(size uint) *unsafe.Pointer {
	if size < PageSize/2 {
		class := uint(math.Pow(2, math.Ceil(math.Log(float64(size))/math.Log(2))))
		page := p.find_divided_class_page(class)
		if page == nil {
			page = p.divide_page(class)
		}
		return p.allocate_single_block(page)
	}
	return p.allocate_multi_block(size)
}

func (p *PageAllocator) allocate_single_block(page *PageHeader) *unsafe.Pointer {
	page.state = DividedPage
	page.blocksAmount--
	page.pointer = unsafe.Pointer(&p.heap[page.adr+page.size])
	page.adr = page.adr + page.size
	return &page.pointer
}

func (p *PageAllocator) allocate_multi_block(size uint) *unsafe.Pointer {
	var sizeFloat float64
	sizeFloat = float64(size)
	needPages := uint(math.Ceil(float64(sizeFloat / PageSize)))
	pagesToSave := []*PageHeader{}
	for _, page := range p.pages {
		if page.state == FreePage {
			pagesToSave = append(pagesToSave, page)
			if len(pagesToSave) == int(needPages) {
				break
			}
		}
	}

	for i, page := range pagesToSave {
		page.state = MultiPage
		page.size = needPages * PageSize
		page.blocksAmount = needPages
		if i == int(needPages)-1 {
			page.blocks = append(page.blocks, BlockHeader{
				adr:  page.adr,
				next: nil,
			})
			break
		}
		page.blocks = append(page.blocks, BlockHeader{
			adr:  page.adr,
			next: unsafe.Pointer(&p.heap[int(page.adr)+(PageSize)]),
		})
	}
	return &pagesToSave[0].pointer
}

func (p *PageAllocator) find_divided_class_page(class uint) *PageHeader {
	for _, page := range p.pages {
		if page.state == DividedPage && page.size == class {
			return page
		}
	}
	return nil
}

func (p *PageAllocator) divide_page(class uint) *PageHeader {
	for _, page := range p.pages {
		if page.state == FreePage {
			page.state = DividedPage
			page.size = class
			page.blocksAmount = PageSize / class

			for i := 0; i < int(page.blocksAmount-1); i++ {
				block := BlockHeader{
					adr:  page.adr + uint(i)*(class),
					next: unsafe.Pointer(&p.heap[page.adr+uint(i+1)*(class)]),
				}
				page.blocks = append(page.blocks, block)
			}
			return page
		}
	}
	return nil
}

func (p *PageAllocator) get_class_by_page(reqPage *PageHeader) (class uint) {
	for _, page := range p.pages {
		if page == reqPage {
			return page.size
		}
	}
	return 0

}

func (p *PageAllocator) mem_realloc(ptr *unsafe.Pointer, size uint) *unsafe.Pointer {
	newPtr := p.mem_alloc(size)
	if newPtr == nil {
		return nil
	}
	p.mem_free(ptr)
	return newPtr
}

func (p *PageAllocator) mem_free(ptr *unsafe.Pointer) {
	for _, page := range p.pages {
		if page.pointer == *ptr {
			if page.state == DividedPage {
				p.free_single_block(page, ptr)
			}
			if page.state == MultiPage {
				p.free_multi_block(page)
			}
		}
	}
}

func (p *PageAllocator) free_multi_block(page *PageHeader) {
	fmt.Println(page.blocksAmount)
	for i := 1; i < int(page.blocksAmount); i++ {
		nextPage := p.get_page_by_adr(page.adr + uint(i)*PageSize)
		nextPage.state = FreePage
		nextPage.blocksAmount = 0
		nextPage.blocks = nil
	}
	page.state = FreePage
	page.blocksAmount = 0
	page.blocks = nil
}

func (p *PageAllocator) get_page_by_adr(adr uint) *PageHeader {
	for _, page := range p.pages {
		if page.adr == adr {
			return page
		}
	}
	return nil
}

func (p *PageAllocator) free_single_block(page *PageHeader, ptr *unsafe.Pointer) {
	for _, block := range page.blocks {
		if unsafe.Pointer(&p.heap[block.adr]) == *ptr {
			page.blocksAmount++
			if page.blocksAmount == PageSize/page.size {
				page.state = FreePage
				page.blocks = nil
				page.size = 0
			}
			break
		}
	}
}

func (p *PageAllocator) mem_dump() {
	fmt.Println("\n|------------Allocator Info------------|")
	fmt.Printf("Memory heap: %v\n", p.heapSize)
	fmt.Printf("Page size: %v\n", PageSize)
	fmt.Printf("Amount of pages: %v\n", p.heapSize/PageSize)
	fmt.Println("|--------------Pages Info--------------|")
	i := 0
	for _, page := range p.pages {
		var statStr string
		switch page.state {
		case FreePage:
			statStr = "Free"
			break
		case DividedPage:
			statStr = "Divided"
			break
		case MultiPage:
			statStr = "Multi"
			break
		}
		fmt.Printf("%v.Status: %v;Address: %v;", i, statStr, page.pointer)
		if page.state == DividedPage {
			fmt.Printf("Class size: %v;Free blocks: %v;\n", p.get_class_by_page(page), page.blocksAmount)
		} else {
			fmt.Printf("\n")
		}
		i++
	}
	fmt.Println("|--------------------------------------|\n")
}

func main() {
	var err error
	var p PageAllocator
	err = p.Init(MemoryHeap)
	if err != nil {
		panic(err)
	}
	fmt.Println("Initialised Page Allocator")
	p.mem_dump()
	fmt.Println("Allocated 8193 bytes")
	ptr8193 := p.mem_alloc(8193)
	p.mem_dump()
	fmt.Println("Allocated 64 and 65 bytes")
	ptr64 := p.mem_alloc(64)
	p.mem_alloc(65)
	p.mem_dump()
	fmt.Println("Reallocated block with 64 bytes to 200")
	ptr200 := p.mem_realloc(ptr64, 200)
	p.mem_dump()
	fmt.Println("Reallocated block with 8193 bytes to 200")
	p.mem_realloc(ptr8193, 200)
	p.mem_dump()
	fmt.Println("Free block with 200")
	p.mem_free(ptr200)
	p.mem_dump()
}
