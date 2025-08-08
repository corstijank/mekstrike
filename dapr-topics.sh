#!/bin/bash

# Dapr Topics Management Script
# Manages Redis streams used by Dapr pub/sub in the mekstrike game

REDIS_POD="redis-master-0"
REDIS_NAMESPACE="redis"
REDIS_PASSWORD="Y5anVVb5C8"

# ANSI color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Dapr topics used by mekstrike
TOPICS=(
    "ai-turn-started"
    "unit-movement-completed"
    "unit-attack-completed"
    "unit-end-phase-completed"
)

redis_exec() {
    kubectl exec $REDIS_POD -n $REDIS_NAMESPACE -- redis-cli -a $REDIS_PASSWORD "$@" 2>/dev/null
}

show_help() {
    echo -e "${BLUE}Dapr Topics Management Script${NC}"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  check, c     - Check all topic stream lengths"
    echo "  list, l      - List recent messages from all topics"
    echo "  clear, cl    - Clear all topic streams"
    echo "  monitor, m   - Monitor topics in real-time (Ctrl+C to stop)"
    echo "  help, h      - Show this help message"
    echo ""
    echo "Topic-specific commands:"
    echo "  check <topic>   - Check specific topic length"
    echo "  list <topic>    - List messages from specific topic"
    echo "  clear <topic>   - Clear specific topic"
    echo ""
    echo "Available topics: ${TOPICS[*]}"
}

check_topics() {
    local topic=$1
    echo -e "${BLUE}=== Topic Stream Lengths ===${NC}"
    
    if [ -n "$topic" ]; then
        # Check specific topic
        local length=$(redis_exec XLEN "$topic")
        if [ "$length" -gt 0 ]; then
            echo -e "${YELLOW}$topic${NC}: $length messages"
        else
            echo -e "${GREEN}$topic${NC}: $length messages"
        fi
    else
        # Check all topics
        for topic in "${TOPICS[@]}"; do
            local length=$(redis_exec XLEN "$topic")
            if [ "$length" -gt 0 ]; then
                echo -e "${YELLOW}$topic${NC}: $length messages"
            else
                echo -e "${GREEN}$topic${NC}: $length messages"
            fi
        done
    fi
    echo ""
}

list_messages() {
    local topic=$1
    
    if [ -n "$topic" ]; then
        # List specific topic
        echo -e "${BLUE}=== Messages from $topic ===${NC}"
        local length=$(redis_exec XLEN "$topic")
        if [ "$length" -gt 0 ]; then
            redis_exec XREAD COUNT 10 STREAMS "$topic" 0 | grep -A1 -E "(data|^$topic)"
        else
            echo "No messages in $topic"
        fi
        echo ""
    else
        # List all topics
        for topic in "${TOPICS[@]}"; do
            echo -e "${BLUE}=== Messages from $topic ===${NC}"
            local length=$(redis_exec XLEN "$topic")
            if [ "$length" -gt 0 ]; then
                echo -e "${YELLOW}Last few messages:${NC}"
                redis_exec XREAD COUNT 3 STREAMS "$topic" 0 | grep -A1 -E "(data|^$topic)" | head -20
            else
                echo "No messages in $topic"
            fi
            echo ""
        done
    fi
}

clear_topics() {
    local topic=$1
    
    if [ -n "$topic" ]; then
        # Clear specific topic
        echo -e "${RED}Clearing topic: $topic${NC}"
        redis_exec DEL "$topic"
        echo -e "${GREEN}✓ Cleared $topic${NC}"
    else
        # Clear all topics
        echo -e "${RED}Clearing all Dapr topics...${NC}"
        for topic in "${TOPICS[@]}"; do
            local length=$(redis_exec XLEN "$topic")
            if [ "$length" -gt 0 ]; then
                redis_exec DEL "$topic"
                echo -e "${GREEN}✓ Cleared $topic ($length messages)${NC}"
            else
                echo -e "${GREEN}✓ $topic (already empty)${NC}"
            fi
        done
    fi
    echo ""
}

monitor_topics() {
    echo -e "${BLUE}=== Monitoring Dapr Topics (Ctrl+C to stop) ===${NC}"
    echo "Watching for new messages..."
    echo ""
    
    # Use MONITOR command to watch Redis activity
    redis_exec MONITOR | grep -E "($(IFS='|'; echo "${TOPICS[*]}"))"
}

# Main script logic
case "$1" in
    "check"|"c")
        check_topics "$2"
        ;;
    "list"|"l")
        list_messages "$2"
        ;;
    "clear"|"cl")
        if [ -n "$2" ]; then
            clear_topics "$2"
        else
            echo -e "${RED}Are you sure you want to clear ALL topics? (y/N)${NC}"
            read -r confirmation
            if [[ "$confirmation" =~ ^[Yy]$ ]]; then
                clear_topics
            else
                echo "Cancelled."
            fi
        fi
        ;;
    "monitor"|"m")
        monitor_topics
        ;;
    "help"|"h"|"")
        show_help
        ;;
    *)
        echo -e "${RED}Unknown command: $1${NC}"
        echo ""
        show_help
        exit 1
        ;;
esac