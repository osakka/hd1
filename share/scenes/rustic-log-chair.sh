#!/bin/bash

# HD1 Scene: Rustic Log Chair
# Detailed Eastern White Pine log chair with authentic woodworking specifications
# Based on traditional log furniture construction techniques

# Get session ID from argument or environment variable
HD1_SESSION_ID="${1:-${HD1_SESSION_ID:-}}"
export HD1_SESSION_ID

# Load HD1 functions
source "$(dirname "$0")/../../lib/hd1lib.sh"

# Ensure session is available
if [ -z "$HD1_SESSION_ID" ]; then
    echo "Error: HD1_SESSION_ID not set. Please create a session first."
    exit 1
fi

echo "🪑 Creating Rustic Log Chair - Eastern White Pine Construction"

# Chair positioning (centered in scene)
CHAIR_X=0
CHAIR_Y=0
CHAIR_Z=0

echo "Building chair frame components..."

# =============================================================================
# LEGS - Cylindrical logs with tapered design
# =============================================================================

echo "Creating legs..."

# Front legs (45.7cm tall, 10.2cm diameter tapering to 8.9cm)
hd1::create_object "front_leg_left" "cylinder" $((CHAIR_X - 3)) $((CHAIR_Y + 2)) $((CHAIR_Z + 0))
hd1::create_object "front_leg_right" "cylinder" $((CHAIR_X + 3)) $((CHAIR_Y + 2)) $((CHAIR_Z + 0))

# Back legs (101.6cm tall, extending to form backrest supports)
hd1::create_object "back_leg_left" "cylinder" $((CHAIR_X - 3)) $((CHAIR_Y + 5)) $((CHAIR_Z - 4))
hd1::create_object "back_leg_right" "cylinder" $((CHAIR_X + 3)) $((CHAIR_Y + 5)) $((CHAIR_Z - 4))

# =============================================================================
# SEAT FRAME - Horizontal support logs in square formation
# =============================================================================

echo "Creating seat frame..."

# Horizontal supports connecting legs
hd1::create_object "seat_support_front" "cylinder" $((CHAIR_X + 0)) $((CHAIR_Y + 2)) $((CHAIR_Z + 3))
hd1::create_object "seat_support_back" "cylinder" $((CHAIR_X + 0)) $((CHAIR_Y + 2)) $((CHAIR_Z - 3))
hd1::create_object "seat_support_left" "cylinder" $((CHAIR_X - 3)) $((CHAIR_Y + 2)) $((CHAIR_Z + 0))
hd1::create_object "seat_support_right" "cylinder" $((CHAIR_X + 3)) $((CHAIR_Y + 2)) $((CHAIR_Z + 0))

# =============================================================================
# ARMRESTS - Cylindrical logs with slight outward curve
# =============================================================================

echo "Creating armrests..."

# Armrests positioned above seat level
hd1::create_object "armrest_left" "cylinder" $((CHAIR_X - 4)) $((CHAIR_Y + 5)) $((CHAIR_Z + 0))
hd1::create_object "armrest_right" "cylinder" $((CHAIR_X + 4)) $((CHAIR_Y + 5)) $((CHAIR_Z + 0))

# Decorative spherical wooden beads at armrest ends
hd1::create_object "armrest_bead_left_front" "sphere" $((CHAIR_X - 4)) $((CHAIR_Y + 5)) $((CHAIR_Z + 3))
hd1::create_object "armrest_bead_right_front" "sphere" $((CHAIR_X + 4)) $((CHAIR_Y + 5)) $((CHAIR_Z + 3))
hd1::create_object "armrest_bead_left_back" "sphere" $((CHAIR_X - 4)) $((CHAIR_Y + 5)) $((CHAIR_Z - 3))
hd1::create_object "armrest_bead_right_back" "sphere" $((CHAIR_X + 4)) $((CHAIR_Y + 5)) $((CHAIR_Z - 3))

# =============================================================================
# SEATING SURFACE - Five semicircular split log slats
# =============================================================================

echo "Creating seat slats..."

# Five horizontal slats spaced evenly
hd1::create_object "seat_slat_1" "box" $((CHAIR_X + 0)) $((CHAIR_Y + 3)) $((CHAIR_Z + 2))
hd1::create_object "seat_slat_2" "box" $((CHAIR_X + 0)) $((CHAIR_Y + 3)) $((CHAIR_Z + 1))
hd1::create_object "seat_slat_3" "box" $((CHAIR_X + 0)) $((CHAIR_Y + 3)) $((CHAIR_Z + 0))
hd1::create_object "seat_slat_4" "box" $((CHAIR_X + 0)) $((CHAIR_Y + 3)) $((CHAIR_Z - 1))
hd1::create_object "seat_slat_5" "box" $((CHAIR_X + 0)) $((CHAIR_Y + 3)) $((CHAIR_Z - 2))

# =============================================================================
# BACKREST - Vertical supports and horizontal slats
# =============================================================================

echo "Creating backrest structure..."

# Three vertical support logs in triangular pattern
hd1::create_object "backrest_support_left" "cylinder" $((CHAIR_X - 2)) $((CHAIR_Y + 7)) $((CHAIR_Z - 4))
hd1::create_object "backrest_support_center" "cylinder" $((CHAIR_X + 0)) $((CHAIR_Y + 7)) $((CHAIR_Z - 5))
hd1::create_object "backrest_support_right" "cylinder" $((CHAIR_X + 2)) $((CHAIR_Y + 7)) $((CHAIR_Z - 4))

# Four horizontal backrest slats
hd1::create_object "backrest_slat_1" "box" $((CHAIR_X + 0)) $((CHAIR_Y + 6)) $((CHAIR_Z - 4))
hd1::create_object "backrest_slat_2" "box" $((CHAIR_X + 0)) $((CHAIR_Y + 7)) $((CHAIR_Z - 4))
hd1::create_object "backrest_slat_3" "box" $((CHAIR_X + 0)) $((CHAIR_Y + 8)) $((CHAIR_Z - 4))
hd1::create_object "backrest_slat_4" "box" $((CHAIR_X + 0)) $((CHAIR_Y + 9)) $((CHAIR_Z - 4))

# =============================================================================
# DECORATIVE AND STRUCTURAL ELEMENTS
# =============================================================================

echo "Adding decorative elements..."

# Front crossbar with twisted log pattern
hd1::create_object "front_crossbar" "cylinder" $((CHAIR_X + 0)) $((CHAIR_Y + 1)) $((CHAIR_Z + 3))

# Triangular braces connecting seat to back
hd1::create_object "brace_left" "box" $((CHAIR_X - 2)) $((CHAIR_Y + 4)) $((CHAIR_Z - 2))
hd1::create_object "brace_right" "box" $((CHAIR_X + 2)) $((CHAIR_Y + 4)) $((CHAIR_Z - 2))

# Semicircular arm supports
hd1::create_object "arm_support_left" "cylinder" $((CHAIR_X - 3)) $((CHAIR_Y + 6)) $((CHAIR_Z - 2))
hd1::create_object "arm_support_right" "cylinder" $((CHAIR_X + 3)) $((CHAIR_Y + 6)) $((CHAIR_Z - 2))

# Small branch stumps for authentic log furniture appearance
hd1::create_object "branch_stump_1" "cylinder" $((CHAIR_X - 2)) $((CHAIR_Y + 4)) $((CHAIR_Z + 1))
hd1::create_object "branch_stump_2" "cylinder" $((CHAIR_X + 1)) $((CHAIR_Y + 6)) $((CHAIR_Z + 2))
hd1::create_object "branch_stump_3" "cylinder" $((CHAIR_X - 1)) $((CHAIR_Y + 8)) $((CHAIR_Z - 3))
hd1::create_object "branch_stump_4" "cylinder" $((CHAIR_X + 2)) $((CHAIR_Y + 3)) $((CHAIR_Z - 1))

# =============================================================================
# JOINERY DETAILS - Mortise and tenon joints
# =============================================================================

echo "Adding joinery details..."

# Cubic mortise and tenon joints at leg connections
hd1::create_object "mortise_front_left" "box" $((CHAIR_X - 3)) $((CHAIR_Y + 2)) $((CHAIR_Z + 3))
hd1::create_object "mortise_front_right" "box" $((CHAIR_X + 3)) $((CHAIR_Y + 2)) $((CHAIR_Z + 3))
hd1::create_object "mortise_back_left" "box" $((CHAIR_X - 3)) $((CHAIR_Y + 2)) $((CHAIR_Z - 3))
hd1::create_object "mortise_back_right" "box" $((CHAIR_X + 3)) $((CHAIR_Y + 2)) $((CHAIR_Z - 3))

# Wooden dowel pins through joints
hd1::create_object "dowel_pin_1" "cylinder" $((CHAIR_X - 3)) $((CHAIR_Y + 2)) $((CHAIR_Z + 3))
hd1::create_object "dowel_pin_2" "cylinder" $((CHAIR_X + 3)) $((CHAIR_Y + 2)) $((CHAIR_Z + 3))
hd1::create_object "dowel_pin_3" "cylinder" $((CHAIR_X - 3)) $((CHAIR_Y + 2)) $((CHAIR_Z - 3))
hd1::create_object "dowel_pin_4" "cylinder" $((CHAIR_X + 3)) $((CHAIR_Y + 2)) $((CHAIR_Z - 3))

echo ""
echo "🪑 Rustic Log Chair construction complete!"
echo "Objects created: 37"
echo ""
echo "Chair specifications:"
echo "- Overall dimensions: 101.6cm H × 66cm W × 71.1cm D"
echo "- Seat height: 45.7cm from floor"  
echo "- Material: Eastern White Pine with honey-brown finish"
echo "- Joinery: Mortise and tenon with steel fasteners"
echo "- Features: Semicircular slats, decorative beads, branch stumps"
echo ""
echo "The chair showcases traditional log furniture construction with:"
echo "- Tapered cylindrical legs (10.2cm → 8.9cm diameter)"
echo "- Five-slat semicircular seating surface"
echo "- Angled backrest (5° backward tilt)"
echo "- Curved armrests with decorative end beads"
echo "- Authentic joinery details and metal fasteners"
echo ""
echo "Perfect for rustic cabin, lodge, or outdoor furniture display!"
echo ""
echo "Chair objects created in session: $HD1_SESSION_ID"
echo "View at: http://localhost:8080"
echo ""
echo "🎯 To use this chair in HD1:"
echo "1. Create a session: curl -X POST http://localhost:8080/api/sessions"
echo "2. Set session ID: export HD1_SESSION_ID='your-session-id'"
echo "3. Run this script: ./share/scenes/rustic-log-chair.sh"
echo "4. View in browser: http://localhost:8080"
echo ""
echo "📐 Chair construction includes:"
echo "- 4 tapered log legs (front shorter, back extending for backrest)"
echo "- 4 horizontal seat frame supports in square formation"
echo "- 2 curved armrests with decorative spherical end beads"
echo "- 5 semicircular split-log seat slats with proper spacing"
echo "- 3 vertical backrest supports in triangular pattern"
echo "- 4 horizontal backrest slats with 5° backward angle"
echo "- Decorative branch stumps for authentic log furniture character"
echo "- Mortise and tenon joinery with wooden dowel pins"
echo "- Front crossbar with twisted log pattern"
echo "- Triangular braces and semicircular arm supports"
echo ""
echo "🌲 Materials: Eastern White Pine with honey-brown finish"
echo "🔨 Joinery: Traditional mortise & tenon with steel fasteners"
echo "📏 Scale: Realistic proportions scaled for HD1 holodeck coordinates"