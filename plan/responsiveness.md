# Responsiveness
## Terminal screen dimensions

widths (cols):
xxxs: 0-24 (standard left 1/3rd split at half width on mon: 11)
xxs: 24 ()
xs: 36 (standard left 1/3rd split mon: 45)
sm: 60 (1/4 width mon: 60)
md: 90 (half width mon: 91)
lg: 120 (3/4 width mon: 121)
xl: 180+ (full width mon: 183)

1.25x screen half width: 106
1.25x screen full width: 219

height (rows):
full_tmux_waybar: 41 - standard
full_tmux_no_waybar: 43
full_no_tmux_or_waybar: 44
half_tmux: 19
half_no_tmux_no_waybar: 21

## Components

### container:

Max width x height container 

xxxs: 0-24 (11)
xxs: 24 ()
xs: 36 (45)
sm: 60 (60)
md: 90 (91)
lg: 120 (121)
xl: 180+ (181)

"╭──────────────────╮  ╭──────────────────╮  ╭──────────────────╮  ╭──────────────────╮"

content should flex to the size of the max width container

### max height container

Manages max


# Terminal notes

  In standard monospace terminal fonts, characters are not
  square—they are typically twice as tall as they are wide (an
  aspect ratio of approximately 1:2 or 0.5:1 for a single
  character cell width to height).

  ### 1. The Terminal Character Cell Ratio

  • 1 column (width) ≈ 0.5 rows (height)
  • Conversely, 1 row (height) ≈ 2 columns (width) in visual
  dimension.

  ### Calculated personal system pixel ratio

  • Cell Width: 6 px
  • Cell Height: 13 px
  ──────
  ### Exact Ratios

  • Row to Column Ratio (Height / Width): 2.167
  • Column to Row Ratio (Width / Height): 0.462

# ratios

4:3
                                                  *0.462*(4/3)
   Box Size Target  | Rows (Height)            | Columns (Width)              
  ------------------|--------------------------|----------------            
   Small            | 6 rows                   | 3.7 columns                
   Small            | 6 rows                   | 4 columns                
   Medium           | 9 rows                   | 5.54 columns               
   Medium           | 9 rows                   | 6 columns               
   Large            | 12 rows                  | 7.39 columns               
   Large            | 12 rows                  | 8 columns               
   Max Fit (Height) | 18 rows                  | 11.09 columns              
   Max Fit (Height) | 18 rows                  | 12 columns              
   Max Fit (Width)  | 27 rows (exceeds height) | 16.63 columns              
                                                                            
 ╭────╮  ╭───────╮ ╭───────────╮                                                                           
 │    │  │       │ │           │
 │    │  │       │ │           │
 ╰────╯  │       │ │           │
         │       │ │           │
         ╰───────╯ │           │                                                                             
                   │           │                                                                             
                   ╰───────────╯                                                                             

3:4                                                                           
   Box Size Target    | Rows (Height)      | Columns (Width)                  
  --------------------|--------------------|--------------------              
   Small              | 4 rows             | 11.6 columns
   Medium             | 6 rows            | 17.336 columns
   Large              | 8 rows            | 23.114666666667 columns
   Max Fit (80x24)    | 12 rows            | 30 columns


╭─────────╮ ╭───────────────╮ ╭─────────────────────╮
│         │ │               │ │                     │
│         │ │               │ │                     │      8*2.167*(3/4)
╰─────────╯ │               │ │                     │
            │               │ │                     │
            ╰───────────────╯ │                     │
                              │                     │
                              ╰─────────────────────╯










