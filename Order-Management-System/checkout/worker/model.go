package worker

import (
	"sync"

	pgxpool "github.com/jackc/pgx/v4/pgxpool"
	"github.com/sofiukl/oms/oms-core/models"
	"github.com/sofiukl/oms/oms-core/utils"
)

// Work - This is WorkRequest model
type Work struct {
	Work   models.CheckoutModel
	Config utils.Config
	Conn   *pgxpool.Pool
	Lock   *sync.RWMutex
}
