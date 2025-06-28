#!/bin/bash

# =========================================================================
# THD Recursion Example: Fractal City Generator
# =========================================================================
#
# Demonstrates recursive programming with:
# - Recursive building generation
# - Mathematical precision
# - Procedural city layouts
# - Infinite complexity from simple rules
# =========================================================================

source "$(dirname "$0")/../thd-shell.sh"

echo "🏙️  THD Fractal City Generator"
echo "=============================="

# City generation parameters
CITY_SIZE=20
MAX_BUILDING_HEIGHT=8
RECURSION_DEPTH=4
BUILDING_DENSITY=0.3

# Create city session
thd::session new "fractal-city-$(date +%s)"

# Initialize city world
echo "🌍 Initializing city world..."
thd::grid $((CITY_SIZE * 2))x$((CITY_SIZE * 2)) 0,0,0 asphalt --name city-ground
thd::camera to $((CITY_SIZE)),20,$((CITY_SIZE))
thd::light sun 0,50,0 yellow 2.0
thd::light ambient 0,0,0 white 0.3

# Recursive building generator
generate_building() {
    local x="$1" z="$2" width="$3" depth="$4" height="$5" recursion_level="$6"
    local building_id="building-${x}-${z}-${recursion_level}"
    
    # Base case: stop recursion
    if [[ $recursion_level -le 0 ]] || [[ $width -lt 1 ]] || [[ $depth -lt 1 ]]; then
        return
    fi
    
    # Create main building
    local y=$(echo "$height / 2" | bc -l)
    thd::cube "$x,$y,$z" gray --scale "$width,$height,$depth" --name "$building_id"
    
    echo "🏢 Building at ($x,$z) - Size: ${width}x${height}x${depth} - Level: $recursion_level"
    
    # Add architectural details
    if [[ $height -gt 3 ]]; then
        # Rooftop structure
        local roof_height=$(echo "$height * 0.2" | bc -l)
        local roof_y=$(echo "$height + $roof_height / 2" | bc -l)
        thd::cube "$x,$roof_y,$z" darkgray --scale "${width}0.8,${roof_height},${depth}0.8" --name "${building_id}-roof"
        
        # Antenna
        if [[ $height -gt 6 ]]; then
            local antenna_height=$(echo "$height * 0.3" | bc -l)
            local antenna_y=$(echo "$height + $antenna_height / 2" | bc -l)
            thd::cube "$x,$antenna_y,$z" red --scale "0.1,$antenna_height,0.1" --name "${building_id}-antenna"
        fi
    fi
    
    # Recursive subdivision - create smaller buildings around
    if [[ $recursion_level -gt 1 ]]; then
        local new_width=$(echo "$width * 0.6" | bc -l)
        local new_depth=$(echo "$depth * 0.6" | bc -l)
        local new_height=$(echo "$height * 0.8" | bc -l)
        local offset=$(echo "$width * 1.5" | bc -l)
        
        # Generate 4 smaller buildings around the main one
        generate_building $((x + offset)) $z "$new_width" "$new_depth" "$new_height" $((recursion_level - 1)) &
        generate_building $((x - offset)) $z "$new_width" "$new_depth" "$new_height" $((recursion_level - 1)) &
        generate_building $x $((z + offset)) "$new_width" "$new_depth" "$new_height" $((recursion_level - 1)) &
        generate_building $x $((z - offset)) "$new_width" "$new_depth" "$new_height" $((recursion_level - 1)) &
    fi
}

# District generator with different architectural styles
generate_district() {
    local center_x="$1" center_z="$2" district_type="$3" size="$4"
    
    echo "🏘️  Generating $district_type district at ($center_x, $center_z)"
    
    case "$district_type" in
        "downtown")
            # Tall, dense buildings
            for i in $(seq 1 5); do
                local x=$((center_x + (RANDOM % size) - size/2))
                local z=$((center_z + (RANDOM % size) - size/2))
                local height=$((6 + RANDOM % 6))
                generate_building "$x" "$z" 3 3 "$height" 2 &
            done
            ;;
        "residential")
            # Smaller, more spread out buildings
            for i in $(seq 1 8); do
                local x=$((center_x + (RANDOM % size) - size/2))
                local z=$((center_z + (RANDOM % size) - size/2))
                local height=$((2 + RANDOM % 3))
                generate_building "$x" "$z" 2 2 "$height" 1 &
            done
            ;;
        "industrial")
            # Wide, low buildings
            for i in $(seq 1 3); do
                local x=$((center_x + (RANDOM % size) - size/2))
                local z=$((center_z + (RANDOM % size) - size/2))
                local width=$((4 + RANDOM % 4))
                local depth=$((4 + RANDOM % 4))
                generate_building "$x" "$z" "$width" "$depth" 3 1 &
            done
            ;;
    esac
}

# Generate street grid
generate_streets() {
    echo "🛣️  Generating street grid..."
    
    # Main avenues (vertical)
    for x in $(seq -$CITY_SIZE 4 $CITY_SIZE); do
        thd::cube "$x,0.1,0" darkgray --scale "1,0.2,$((CITY_SIZE * 2))" --name "avenue-$x"
    done
    
    # Main streets (horizontal) 
    for z in $(seq -$CITY_SIZE 4 $CITY_SIZE); do
        thd::cube "0,0.1,$z" darkgray --scale "$((CITY_SIZE * 2)),0.2,1" --name "street-$z"
    done
    
    # Intersections with traffic lights
    for x in $(seq -$CITY_SIZE 8 $CITY_SIZE); do
        for z in $(seq -$CITY_SIZE 8 $CITY_SIZE); do
            thd::cube "$x,2,$z" yellow --scale "0.3,3,0.3" --name "light-$x-$z"
            thd::rotate "light-$x-$z" y 1.0  # Blinking effect
        done
    done
}

# Populate with vehicles and life
add_city_life() {
    echo "🚗 Adding city life..."
    
    # Cars on streets
    for i in $(seq 1 20); do
        local x=$(( (RANDOM % (CITY_SIZE * 2)) - CITY_SIZE ))
        local z=$(( (RANDOM % (CITY_SIZE * 2)) - CITY_SIZE ))
        local colors=("red" "blue" "white" "black" "silver")
        local color="${colors[$((RANDOM % ${#colors[@]}))]}"
        
        thd::cube "$x,0.5,$z" "$color" --scale "2,1,1" --name "car-$i"
        
        # Random movement
        thd::after "$((RANDOM % 10))s" "thd::move car-$i by $((RANDOM % 6 - 3)),0,$((RANDOM % 6 - 3)) in 5s"
    done
    
    # People (small spheres)
    for i in $(seq 1 50); do
        local x=$(( (RANDOM % (CITY_SIZE * 2)) - CITY_SIZE ))
        local z=$(( (RANDOM % (CITY_SIZE * 2)) - CITY_SIZE ))
        
        thd::sphere "$x,0.3,$z" random --scale "0.3,0.3,0.3" --name "person-$i"
        
        # Random walking
        thd::every "$((2 + RANDOM % 5))s" "thd::move person-$i by $((RANDOM % 3 - 1)),0,$((RANDOM % 3 - 1)) in 1s" &
    done
}

# Cinematic city tour
city_tour() {
    echo "🎬 Starting cinematic city tour..."
    
    # Tour waypoints
    local waypoints=(
        "-20,5,-20"    # Approach the city
        "0,10,0"       # City center overview  
        "10,3,10"      # Street level
        "0,25,0"       # High aerial view
        "20,5,20"      # Departure
    )
    
    for waypoint in "${waypoints[@]}"; do
        echo "📍 Moving to waypoint: $waypoint"
        thd::camera to "$waypoint" in 4s
        thd::wait 5s
    done
    
    # Final orbital shot
    echo "🌍 Final orbital sequence"
    thd::camera orbit 0,0,0 --radius 30 --speed 0.5
}

# Main city generation sequence
echo "🚀 Starting fractal city generation..."

# Phase 1: Infrastructure
echo "📍 Phase 1: Infrastructure"
generate_streets
thd::wait 2s

# Phase 2: Districts
echo "📍 Phase 2: District Generation"
thd::parallel {
    generate_district -8 -8 "downtown" 8
    generate_district 8 -8 "residential" 12  
    generate_district -8 8 "industrial" 10
    generate_district 8 8 "residential" 12
}

echo "⏳ Waiting for buildings to complete..."
thd::wait 10s

# Phase 3: Central landmark (recursive skyscraper)
echo "📍 Phase 3: Central Landmark"
echo "🏗️  Building recursive skyscraper..."
generate_building 0 0 4 4 12 $RECURSION_DEPTH

thd::wait 5s

# Phase 4: City life
echo "📍 Phase 4: Adding Life to the City"
add_city_life

thd::wait 3s

# Phase 5: Cinematic tour
echo "📍 Phase 5: Cinematic Tour"
city_tour

# Final statistics
echo ""
echo "🏙️  FRACTAL CITY COMPLETE!"
echo "=========================="
echo "City Size: ${CITY_SIZE}x${CITY_SIZE}"
echo "Recursion Depth: $RECURSION_DEPTH" 
echo "Building Count: ~$(($(seq 1 5 | wc -l) * 4))+ recursive structures"
echo "Street Grid: Complete with intersections"
echo "Population: 70+ entities (cars + people)"
echo "Tour Duration: ~30 seconds"
echo ""
echo "🎭 Ready for your next fractal creation!"