#!/bin/bash

# =========================================================================
# THD Game Example: 3D Snake Game
# =========================================================================
#
# Demonstrates game programming with:
# - Real-time player control
# - Game state management
# - Collision detection
# - Score tracking and growth mechanics
# =========================================================================

source "$(dirname "$0")/../thd-shell.sh"

echo "🐍 THD Game: 3D Snake"
echo "===================="

# Game configuration
SNAKE_LENGTH=3
SNAKE_DIRECTION="right"
FOOD_COUNT=5
SCORE=0
GAME_SPEED=0.5

# Create game session
thd::session new "snake-game-$(date +%s)"

# Initialize game world
echo "🎮 Initializing game world..."
thd::grid 20x20 0,0,0 darkgray --name game-board
thd::camera to 10,15,10
thd::light game 0,10,0 white 1.5

# Create snake body
echo "🐍 Creating snake..."
SNAKE_BODY=()
for i in $(seq 0 $((SNAKE_LENGTH-1))); do
    local x=$((5 + i))
    thd::cube "$x,1,10" green --name "snake-$i"
    SNAKE_BODY+=("snake-$i")
done

# Create initial food
echo "🍎 Spawning food..."
spawn_food() {
    local food_id="$1"
    local x=$((RANDOM % 20))
    local z=$((RANDOM % 20))
    thd::sphere "$x,1,$z" red --name "food-$food_id"
    thd::rotate "food-$food_id" y 2.0  # Spinning food
}

for i in $(seq 1 $FOOD_COUNT); do
    spawn_food "$i"
done

# Game controls simulation (in real implementation, would bind to actual keys)
echo "🎮 Game Controls Ready (WASD to move)"

# Main game loop
echo "🚀 Starting game loop..."
GAME_RUNNING=true

# Simulate game input (in real game, this would be actual key input)
simulate_input() {
    local directions=("up" "down" "left" "right")
    local new_direction="${directions[$((RANDOM % 4))]}"
    
    # Prevent reverse direction
    case "$SNAKE_DIRECTION-$new_direction" in
        "up-down"|"down-up"|"left-right"|"right-left")
            return  # Invalid move
            ;;
        *)
            SNAKE_DIRECTION="$new_direction"
            ;;
    esac
}

# Move snake in current direction
move_snake() {
    # Get head position
    local head_pos=$(thd_get_position "${SNAKE_BODY[0]}")
    local x y z
    IFS=',' read -r x y z <<< "$head_pos"
    
    # Calculate new head position
    case "$SNAKE_DIRECTION" in
        "up")    z=$((z - 1)) ;;
        "down")  z=$((z + 1)) ;;
        "left")  x=$((x - 1)) ;;
        "right") x=$((x + 1)) ;;
    esac
    
    # Check boundaries
    if [[ $x -lt 0 || $x -ge 20 || $z -lt 0 || $z -ge 20 ]]; then
        echo "💀 Game Over: Hit boundary!"
        GAME_RUNNING=false
        return
    fi
    
    # Check self-collision
    for segment in "${SNAKE_BODY[@]:1}"; do
        local seg_pos=$(thd_get_position "$segment")
        if [[ "$seg_pos" == "$x,$y,$z" ]]; then
            echo "💀 Game Over: Hit self!"
            GAME_RUNNING=false
            return
        fi
    done
    
    # Move body segments (from tail to head)
    for i in $(seq $((${#SNAKE_BODY[@]} - 1)) -1 1); do
        local prev_pos=$(thd_get_position "${SNAKE_BODY[$((i-1))]}")
        thd::move "${SNAKE_BODY[$i]}" to "$prev_pos"
    done
    
    # Move head to new position
    thd::move "${SNAKE_BODY[0]}" to "$x,$y,$z"
    
    # Check food collision
    check_food_collision "$x,$y,$z"
}

# Check if snake ate food
check_food_collision() {
    local head_pos="$1"
    
    for i in $(seq 1 $FOOD_COUNT); do
        local food_pos=$(thd_get_position "food-$i")
        if [[ "$food_pos" == "$head_pos" ]]; then
            # Food eaten!
            echo "🍎 Food eaten! Score: $((++SCORE))"
            
            # Remove eaten food
            thd::delete "food-$i"
            
            # Grow snake
            grow_snake
            
            # Spawn new food
            spawn_food "$i"
            
            # Increase speed slightly
            GAME_SPEED=$(echo "$GAME_SPEED * 0.95" | bc -l)
            
            break
        fi
    done
}

# Grow snake by adding segment
grow_snake() {
    local tail_pos=$(thd_get_position "${SNAKE_BODY[-1]}")
    local new_segment="snake-${#SNAKE_BODY[@]}"
    
    thd::cube "$tail_pos" green --name "$new_segment"
    SNAKE_BODY+=("$new_segment")
    
    echo "🐍 Snake grew! Length: ${#SNAKE_BODY[@]}"
}

# Placeholder for getting object position (would integrate with THD API)
thd_get_position() {
    local object="$1"
    # In real implementation, would query THD server for object position
    echo "5,1,10"  # Placeholder
}

# Game loop with timing
game_loop() {
    local frame=0
    
    while [[ "$GAME_RUNNING" == "true" ]]; do
        echo "🎮 Frame $((++frame)) - Score: $SCORE - Length: ${#SNAKE_BODY[@]}"
        
        # Simulate random input occasionally
        if [[ $((frame % 5)) -eq 0 ]]; then
            simulate_input
        fi
        
        # Move snake
        move_snake
        
        # Wait for next frame
        thd::wait "$GAME_SPEED"s
        
        # Auto-quit after demo period
        if [[ $frame -gt 50 ]]; then
            echo "🏁 Demo complete!"
            break
        fi
    done
}

# Start the game
echo "🎮 Game starting in 3 seconds..."
thd::wait 3s

game_loop

# Game over effects
echo "🎭 Game Over Sequence"
thd::effect fade out 2s
thd::camera orbit 10,1,10 --radius 15

# Final score display
echo ""
echo "🐍 GAME OVER 🐍"
echo "==============="
echo "Final Score: $SCORE"
echo "Snake Length: ${#SNAKE_BODY[@]}"
echo "Frames Played: $frame"
echo ""
echo "Ready for another round! 🎮"