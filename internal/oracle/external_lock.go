package oracle

// PSEUDOCODE

import (
	"github.com/ethereum/go-ethereum/common"
)

// SINGLETON
func GetExternalLockServer() *ExternalLockServer {
	externalOnce.Do(func() {
		//ako server nije pokrenut, pokrene se
	})
	return externalLock
}

func (s *ExternalLockServer) bookSymbolWithTx(symbol string) (string, bool) {
	lock, reason, err := bookSymbolWithTx_API_poziv_ka_serveru(symbol) // po prijemu api poziva server automatski bukira simbol ako nije bukiran(postavlja lock na njega)
	// vraca:	 lock - da li imamo lock na tom simbolu
	//			 reason (string) - razlog zakljucavanja (ako je trenutno poznat):
	//						1) hash transakcije ciji se status finalizacije ceka
	//						2) "currently UNKNOWN" ukoliko je API poslat bas u malom periodu nakon lockovanja simbola na serveru
	// 						 	pre nego sto server dobije hash vrednost potpisane transakcije koja je izazvala lock i ciju finalizaciju cekamo
	//			 err - da li je API poziv uspeo (mozda server ne radi)

	if err != nil { // ako server ne radi (ne odaziva se)
		return "", false // onda nista nije zakljucano i aplikacija nastavlja da radi bez spoljnog servera
	}

	if lock == false { // ako simbol nije bukiran
		return "", false // slobodno nastavi izvrsavanje programa
	}

	if lock == true { //ako je simbol vec bukiran, vracamo razlog bukiranja i lock=true (pa se setPrice za dati simbol preskace u ovih 1min)
		return reason, true
	}

	return "", false
}

func (s *ExternalLockServer) ProvideTxHashToServer(symbol string, txHash common.Hash) {
	POST_API_poziv_ka_serveru(symbol, txHash) //saljem serveru hash transakcije koji mu dugujem za zakljucan simbol
}

/*
Server je pokrenut negde drugde i u svojoj implementaciji podrzava:
	- 2 APIja
	- mapu ciji je kljuc simbol, a vrednost par <lock, txHash>

1. bookSymbolWithTx_API_poziv_ka_serveru(symbol)
	ako simbol nije zakljucan:
		- automatski ga zakljucava
		- vraca -> lock=false, reason = ""

	ako je simbol zakljucan:
		- vraca -> lock=true, reason = sadrzaj iz mape(txHash ili currently UNKNOWN)


2. POST_API_poziv_ka_serveru(symbol, txHash)
	- preko njega server dobija txHash transakcije koja je lokovala simbol i ciju finalizaciju server treba da prati
	- u mapu za simbol upisujemo -> txHash   		//da bismo imali razlog zakljucavanja koji vracamo ostalim instancama koji pozivaju 1. API
	- automatski pustamo gorutinu(za txHash) koja po finalizaciji statusa otkljucava odgovarajuci simbol u mapi

*/
