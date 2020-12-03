# page_memory_allocator
In this realisation of Page Memory Allocator we got 3 structures.
```
type PageAllocator struct {
	heapSize     uint
	heap         []byte
	pages        []*PageHeader
}
```
PageAllocator contains of 3 fields:
* heapSize - is size of memory heap we work with in bytes,in example it's equals **32768 bytes**;
* heap - is actualy part of memory we work with;
* pages - is slice of PageHeaders that contains all info about pages on which divided memory heap;
```
type PageHeader struct {
	state        uint
	size         uint
	adr          uint
	pointer      unsafe.Pointer
	blocksAmount uint
	blocks       []BlockHeader
}
```
PageHeader contains of 6 fields:
* state  - iint value that discribes state in which page is in the moment(**Free,Divided or Multi**);
* size - if page is in Divided state it field dsscribes class to which it represent;
* adr - is adress of next free block in this page;
* pointer  - pointer on this page in memory heap;
* blocksAmount - if page is Divided it's shows how much blocks are free to use,if it's Multi it's dicribes amount of pages in multipage;
* blocks - contains slice of BlockHeader;


```
type BlockHeader struct {
	adr  uint
	next unsafe.Pointer
}
```
BlockHeader contains of 2 fields:
* adr  - is adress of this block in memory heap;
* next- pointer to next block if it is so;

Based on this structures were created methods **mem_alloc,mem_realloc,mem_free** and **mem_dump**;
Example of working with which is situated bellow.

* **mem_dump**
```
Initialised Page Allocator

|------------Allocator Info------------|
Memory heap: 32768
Page size: 4096
Amount of pages: 8
|--------------Pages Info--------------|
0.Status: Free;Address: 0x7fc696cfc000;
1.Status: Free;Address: 0x7fc696cfd000;
2.Status: Free;Address: 0x7fc696cfe000;
3.Status: Free;Address: 0x7fc696cff000;
4.Status: Free;Address: 0x7fc696d00000;
5.Status: Free;Address: 0x7fc696d01000;
6.Status: Free;Address: 0x7fc696d02000;
7.Status: Free;Address: 0x7fc696d03000;
|--------------------------------------|
```
* **mem_alloc**
```
Allocated 8193 bytes

|------------Allocator Info------------|
Memory heap: 32768
Page size: 4096
Amount of pages: 8
|--------------Pages Info--------------|
0.Status: Multi;Address: 0x7fc696cfc000;
1.Status: Multi;Address: 0x7fc696cfd000;
2.Status: Multi;Address: 0x7fc696cfe000;
3.Status: Free;Address: 0x7fc696cff000;
4.Status: Free;Address: 0x7fc696d00000;
5.Status: Free;Address: 0x7fc696d01000;
6.Status: Free;Address: 0x7fc696d02000;
7.Status: Free;Address: 0x7fc696d03000;
|--------------------------------------|
```
```
Allocated 64 and 65 bytes

|------------Allocator Info------------|
Memory heap: 32768
Page size: 4096
Amount of pages: 8
|--------------Pages Info--------------|
0.Status: Multi;Address: 0x7fc696cfc000;
1.Status: Multi;Address: 0x7fc696cfd000;
2.Status: Multi;Address: 0x7fc696cfe000;
3.Status: Divided;Address: 0x7fc696cff040;Class size: 64;Free blocks: 63;
4.Status: Divided;Address: 0x7fc696d00080;Class size: 128;Free blocks: 31;
5.Status: Free;Address: 0x7fc696d01000;
6.Status: Free;Address: 0x7fc696d02000;
7.Status: Free;Address: 0x7fc696d03000;
|--------------------------------------|

```
* **mem_realloc**
```
Reallocated block with 64 bytes to 200

|------------Allocator Info------------|
Memory heap: 32768
Page size: 4096
Amount of pages: 8
|--------------Pages Info--------------|
0.Status: Multi;Address: 0x7fc696cfc000;
1.Status: Multi;Address: 0x7fc696cfd000;
2.Status: Multi;Address: 0x7fc696cfe000;
3.Status: Free;Address: 0x7fc696cff040;
4.Status: Divided;Address: 0x7fc696d00080;Class size: 128;Free blocks: 31;
5.Status: Divided;Address: 0x7fc696d01100;Class size: 256;Free blocks: 15;
6.Status: Free;Address: 0x7fc696d02000;
7.Status: Free;Address: 0x7fc696d03000;
|--------------------------------------|
```
```

Reallocated block with 8193 bytes to 200
3

|------------Allocator Info------------|
Memory heap: 32768
Page size: 4096
Amount of pages: 8
|--------------Pages Info--------------|
0.Status: Free;Address: 0x7fc696cfc000;
1.Status: Free;Address: 0x7fc696cfd000;
2.Status: Free;Address: 0x7fc696cfe000;
3.Status: Free;Address: 0x7fc696cff040;
4.Status: Divided;Address: 0x7fc696d00080;Class size: 128;Free blocks: 31;
5.Status: Divided;Address: 0x7fc696d01200;Class size: 256;Free blocks: 14;
6.Status: Free;Address: 0x7fc696d02000;
7.Status: Free;Address: 0x7fc696d03000;
|--------------------------------------|

```
* **mem_free**
```
Free block with 200

|------------Allocator Info------------|
Memory heap: 32768
Page size: 4096
Amount of pages: 8
|--------------Pages Info--------------|
0.Status: Free;Address: 0x7fc696cfc000;
1.Status: Free;Address: 0x7fc696cfd000;
2.Status: Free;Address: 0x7fc696cfe000;
3.Status: Free;Address: 0x7fc696cff040;
4.Status: Divided;Address: 0x7fc696d00080;Class size: 128;Free blocks: 31;
5.Status: Divided;Address: 0x7fc696d01200;Class size: 256;Free blocks: 15;
6.Status: Free;Address: 0x7fc696d02000;
7.Status: Free;Address: 0x7fc696d03000;
|--------------------------------------|

```

