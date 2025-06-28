#!/bin/bash

# =========================================================================
# THD Movie Example: Epic Battle Scene
# =========================================================================
#
# Demonstrates cinematic 3D programming with:
# - Timing and choreography
# - Parallel execution
# - Camera cinematography
# - Effects and drama
# =========================================================================

source "$(dirname "$0")/../thd-shell.sh"

echo "🎬 THD Movie: Epic Battle Scene"
echo "==============================="

# Create new session for the movie
thd::session new "epic-battle-$(date +%s)"

# Scene 1: Setup the battlefield
echo "🎭 Scene 1: Setting the Stage"
thd::grid 50x50 0,0,0 darkgreen --name battlefield
thd::light sun 0,20,0 yellow 2.0
thd::camera to 0,15,30 in 3s

thd::wait 2s

# Scene 2: Armies assemble
echo "⚔️  Scene 2: Armies Assemble"

# Red army (left side)
thd::parallel {
    for i in {1..20}; do
        local x=$((i % 5 - 10))
        local z=$((i / 5 - 10))
        thd::cube "$x,1,$z" red --name "red-$i" &
        thd::wait 0.1s
    done
    wait
}

# Blue army (right side)  
thd::parallel {
    for i in {1..20}; do
        local x=$((i % 5 + 6))
        local z=$((i / 5 - 10))
        thd::sphere "$x,1,$z" blue --name "blue-$i" &
        thd::wait 0.1s
    done
    wait
}

thd::wait 3s

# Scene 3: Dramatic pause with camera work
echo "📷 Scene 3: Dramatic Buildup"
thd::camera fly-through -15,5,0 to 15,5,0 in 4s
thd::effect fade in 1s

thd::wait 5s

# Scene 4: The charge begins
echo "🏃 Scene 4: The Charge!"

# Red army charges
thd::parallel {
    thd::move red-* to +8,1,0 in 3s
    thd::effect particle dust -5,0,0 200
    
    # Blue army charges
    thd::move blue-* to -8,1,0 in 3s  
    thd::effect particle dust 5,0,0 200
    
    # Dramatic camera follow
    thd::camera orbit 0,1,0 --radius 20 --speed 2.0
}

thd::wait 4s

# Scene 5: Epic clash with effects
echo "💥 Scene 5: Epic Clash!"

# Collision explosions
thd::every 0.5s "thd::effect explosion random --size 2.0" &
explosion_pid=$!

# Chaos - objects flying everywhere
thd::parallel {
    for obj in red-* blue-*; do
        thd::rotate "$obj" y 5.0 &
        thd::scale "$obj" 0.5 in 1s &
        thd::after "$((RANDOM % 3))s" "thd::move $obj to random" &
    done
}

thd::wait 5s

# Scene 6: Climactic finale
echo "🌟 Scene 6: Finale"

# Stop explosions
kill $explosion_pid 2>/dev/null

# Victory celebration
thd::effect fade out 2s
thd::light victory 0,10,0 gold 3.0
thd::camera to 0,25,25 in 3s

# Spiraling victory dance
thd::spiral all around 0,5,0 in 10s

thd::wait 8s

# Credits
echo "🎭 Scene 7: Credits"
thd::effect fade to black 3s
thd::wait 2s

echo ""
echo "🎬 Epic Battle Scene Complete!"
echo "================================"
echo "✅ Cinematic timing achieved"
echo "✅ Parallel execution mastered"  
echo "✅ Camera choreography perfected"
echo "✅ Effects and drama delivered"
echo ""
echo "Ready for your next THD movie! 🎭"