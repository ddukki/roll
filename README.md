# `roll` Command Line Tool

The `roll` command line tool provides a simple utility to roll dice through the command line. The tool runs very simple dice roll expressions that are common in D&D game settings, but is also useful for general dice rolling.

The code was written to learn more about parsing syntax trees and also to facilitate quick dice rolls during games of chance and role-playing.

## General Dice Rolling

To roll a dice with `n` faces, the `dn` command can be used. For example, to roll a 20-sided die, run the following command:

```sh
roll d20
```

## Roll Expressions

To add literal values to the rolls, simple math expressions can be used to augment the rolls. For example, to roll a 6-sided die and add 2 to the total, run the following command:

```sh
roll d6 + 2
```

Multiple dice rolls can also be added to the expression:

```sh
# will roll a 20-sided die, a 6-sided die and subtract 5 from the sum
roll d20 + d6 - 5 
```

Currently the following operations are supported:

| Op | Description |
|---|---|
| + | Add |
| - | Subtract |
| * | Multiply |
| / | Divide |

## Multi-Dice Rolls

To roll multiple dice with the same number of faces, a simplified expression can be used instead of using multiple math operations. For example, the following shows different ways to roll four 8-sided dice:

```sh
# with math expressions
roll d8 + d8 + d8 + d8 

# with a simplified dice roll expression
roll 4d8
```

## Complex Dice Expressions

In some cases, the max or min value of a set of rolls may be desired, instead of the sum of all dice. Or, it may be desired that those values be removed from the sum. These operations can be done by using `h`, `H`, `l`, and `L`.

| Expression | Description |
|---|---|
| h | Keep highest (max) of a multi-dice roll |
| H | Drop highest from a multi-dice roll sum |
| l | Keep lowest (min) of a multi-dice roll |
| L | Drop lowest from a multi-dice roll sum |

## Possible Future Features

- Parentheses
- Exponentiation
- Modulus
- Detecting "Nat" values
- Verbose/Sparse output options
- Dice roll history/log
- Non-pipped Dice
- Mapping values to a lookup table
- Help section with `--help` flag
- Piping the output to another project for further utility