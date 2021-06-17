# Century golem edition



## Models

These are the models which will be used to represent our game. 

### Coins
There are a couple of contexts by which we may want to represent coins in this game (player posession, availabilty in game)
in both cases, we will represent these with the same object:

```
{
  copper: int,
  silver: int
}
```

### Resources
There are many contexts in the game in which we need to represent the number of gems (inputs, costs, outputs, stacked on resource cards, etc.)

These will all be represented with the same object
```
{
  yellow: int,
  green: int,
  blue: int,
  pink: int
}
```

### Resource cards
Resource cards are the cards that make up a players hand. These cards are used to generate or convert gems and will be represented as follows:

input and output are represented by the resource object defined above. The id is the string representation of the card. This is a series of characters that uniquely identify a card.
_with any id that includes a resource value serialized, you can omit any 0 gem values_

example resource_ids: [`2y>1b`, `>3u`, `>3y`]
The above cards are: 2 yellow converted into 1 blue; get 3 upgrades; and get 3 yellows

These can be interpreted as cards by doing the following: the `>` character is required for all resource id's and splits the inputs and outputs. Inputs and outputs are defined by a series of number, letter pairs that identify the count and color of gems input or output by the card. Additionally, the output string can include a number of upgrades represented by `{number}u` where number is the count of upgrades the card provides (either 2 or 3)
```
{
  id: resource_id,
  input: resource_object,
  output: resource_object,
  upgrades: int
}
```

### Golem cards
Golem cards are the cards by which players score the majority of their points. These cards are represented by their gem cost and the variable point values as well as a string which allows for easy communication about a card (in an smaller format)

The golem_id is defined by a series of number letter pairs followed by the point value they are worth.
example golem_ids: [`2y2g2b2p20`, `5p20`] 
the above cards can be interpreted as follows: a golem that costs 2 of each color resource and is worth 20 points; A golem that costs 5 pinks and is worth 20 points.
 
_with any id that includes a resource value serialized, you can omit any 0 gem values_

```
{
  id: golem_id,
  cost: resource_object,
  points: int
}
```

### Game (serialized by as the gamestate)
This is the representation of the game. This is how all information about a specific game will be stored / represented by the backend
``` {
  turn_tracker: player_id,
  player_ids: [string_id, ... , string_id]
  available_golems: [{}, {}, {} {} {}],
  available_resources: [{}, {}, {}, {}, {}, {}],

  available_coins: {
    copper: int,
    silver: int
  },



  golem_deck: [{}, ... , {}]
  resource_deck: [{}, ... , {}]
  

  player_id: {
      coins: {
        copper: int,
        silver: int
      },
      golems: [
        {
          id: string_id,
          points: int,
          cost: {
              yellow: int,
              green: int,
              blue: int,
              pink: int
          }
        }, {
          ...
        }
      ],
      resources: [
        {
          id: string_id
          cost: {
            yellow: int,
            green: int,
            blue: int,
            pink: int
          }
          output: {
            yellow: int,
            green: int,
            blue: int,
            pink: int,
            upgrades: int
          }
        }, {
          ...
        }
      ]
      discard_pile: [
        {
          ...
        }, {
          ...
        }
      ]
      gems: {
        yellow: int,
        green: int,
        blue: int,
        pink: int
      }
  },


  player_id:  {
    ...
  }
}
```

The api will return a slightly modified version of this json blob as to hide some information about the opponents' hands and hidden information such as deck order

## Actions

### Create a game lobby
This is the endpoint for creating a game lobby

### Join a game lobby

### Start a game

### Get game state

### Perform Action


