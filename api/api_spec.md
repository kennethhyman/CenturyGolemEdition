# Century golem edition



## Models

These are the models which will be used to represent our game. 

### Game (serialized by as the gamestate)
This is the representation of the game. This is how all information about a specific game will be stored / represented by the backend
```
{
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


