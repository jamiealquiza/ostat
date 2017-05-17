package main

// Stat holds information fetched.
type Stat map[string]struct {
	General struct {
		Uptime int64
		CPU    struct {
			Model string
			Cores int
		}
		Load struct {
			Short float64
			Mid   float64
			Long  float64
			Procs uint16
		}
		Mem struct {
			Total     uint64
			Free      uint64
			Used      uint64
			Usedp     uint64
			Shared    uint64
			Buffer    uint64
			Swaptotal uint64
			Swapfree  uint64
		}
	}
	Storage map[string]Storage
}

// Storage holds storage-specific
// information.
type Storage struct {
	Free        uint64
	Inodesfree  uint64
	Inodestotal uint64
	Inodesused  uint64
	Total       uint64
	Type        string
	Used        uint64
	Usedp       uint64
}
