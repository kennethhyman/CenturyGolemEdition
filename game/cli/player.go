package cli


//import (
//  pb "github.com/kennethhyman/CenturyGolemEdition/grpc"
//)

type Player struct {
  hand []GemCard
  golems []GolemCard
  discards []GemCard
  gems GemValue
  gold_coins int8
  silver_coins int8
}

func (p Player) String() string {
  // name : coins | golems 
  // gems [discard_count]
  // gem cards | 
  return ""
}
