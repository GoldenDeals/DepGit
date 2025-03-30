package git

import "github.com/GoldenDeals/DepGit/internal/stroage"

type Server struct {
	config  Config
	stroage stroage.Stroage
}
