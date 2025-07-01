# 🌊 **Flow – Protect Your Attention**

*A session-based ritual for deep work, designed in the spirit of Zenta.*

---

## 🧭 Purpose

Flow is a mindful terminal tool that helps developers:

* Create **sacred boundaries** around deep work.
* Maintain **one focused session** at a time.
* Work with **presence**, not productivity pressure.

It does **one thing**: create a quiet, bounded space for conscious work — with full awareness.

---

## 🌱 Design Principles

| Principle         | Meaning in Flow                          |
| ----------------- | ---------------------------------------- |
| **Mindful**       | Ritual approach, gentle awareness.       |
| **Single-focus** | One session at a time. No multitasking. |
| **Background**    | Runs quietly, doesn't interfere.        |
| **Unix-like**     | Clean commands, composable design.      |
| **Zenta spirit** | Serves consciousness, not metrics.       |

---

## 🧱 Core Functionality

### ✅ Session Commands

```bash
flow start --tag "deep work"    # Begin focused session
flow status                     # Check current work
flow pause                      # Mindful break
flow resume                     # Continue working  
flow end                        # Complete session
```

### ✅ What It Does

* Creates a **mindful container** for deep work
* **Enforces single-tasking** - prevents multiple sessions
* **Runs in background** - terminal stays clean for work
* **Gentle awareness** - shows current focus when asked
* **Natural rhythms** - supports pause/resume cycles

---

## 🖼️ Terminal Experience

### Starting Work
```
$ flow start --tag "writing documentation"

🌊 Starting deep work: writing documentation

   Clear your mind
   Focus on what matters
   Let distractions pass

Session active in background.
Use 'flow status' to check, 'flow end' to complete.
```

### Checking Status
```
$ flow status
🌊 Deep work: writing documentation
Active for 1h 23m.
```

### Session Management
```
$ flow pause
⏸️  Paused: writing documentation
Worked for 1h 23m. Use 'flow resume' when ready.

$ flow resume
🌊 Resumed: writing documentation
Continue your deep work.

$ flow end
✨ Session complete: writing documentation
Total focus time: 2h 15m

Carry this focus forward.
```

---

## 🎮 Command Reference

| Command | Purpose | Options | Example |
|---------|---------|---------|---------|
| `start` | Begin session | `--tag "description"` | `flow start --tag "refactoring"` |
| `status` | Check current work | None | `flow status` |
| `pause` | Take mindful break | None | `flow pause` |
| `resume` | Continue session | None | `flow resume` |
| `end` | Complete work | None | `flow end` |

---

## 🧘 Session Enforcement

### One Thing at a Time
```bash
$ flow start --tag "task A"
🌊 Starting deep work: task A

$ flow start --tag "task B"  
🌊 Already in deep work: task A
Working for 15m. Use 'flow end' to complete.

One thing at a time.
```

### Natural Work Cycles
- **Start**: Mindful intention setting
- **Work**: Clean terminal, background awareness
- **Pause**: Gentle break without judgment
- **Resume**: Seamless continuation
- **End**: Completion acknowledgment

---

## 🔐 Privacy & Storage

> All data stays local. No tracking. No cloud.

* Session state stored in `~/.flow-session`
* **No productivity logs** - only current session awareness
* **No metrics collection** - focus on presence, not performance
* **Temporary storage** - session file removed on completion

---

## 🧘 Integration with [Zenta](https://github.com/e6a5/zenta)

Natural composition for mindful work cycles:

```bash
flow start --tag "deep coding"     # Enter focus
# ... work happens ...
flow pause && zenta now             # Mindful break
flow resume                         # Continue work
flow end && zenta                   # Breathe after completion
```

Flow creates boundaries; Zenta provides breath. Perfect companions.

---

## 🛠 Architecture

### Single File Design
```
main.go - Complete session management
├── Command parsing (start, status, pause, resume, end)
├── Session state management (JSON file)
├── Mindful display formatting
└── Unix-like exit codes
```

### Session State
```json
{
  "tag": "writing documentation",
  "start_time": "2024-01-15T10:00:00Z",
  "paused_at": "2024-01-15T11:30:00Z",
  "is_paused": true,
  "total_paused": "15m30s"
}
```

---

## ✅ Implementation Status

| Feature | Status | Description |
|---------|--------|-------------|
| Session start | ✅ | Mindful intention setting |
| Background operation | ✅ | Quiet, non-intrusive |
| Status checking | ✅ | On-demand awareness |
| Pause/resume | ✅ | Natural work rhythms |
| Session completion | ✅ | Gentle acknowledgment |
| Single session enforcement | ✅ | One thing at a time |
| Zenta composition | ✅ | Shell operator friendly |

---

## 🌊 Philosophy

### Real vs Fake Focus Tools

**✅ Real mindful focus (Flow's way):**
- Creates sacred space for work
- Respects natural attention rhythms  
- No measurement or optimization
- Serves consciousness, not productivity

**❌ Fake productivity focus:**
- Tracks and measures work output
- Gamifies attention with scores
- Optimizes performance metrics
- Makes focus about achievement

---

## 💬 Core Insight

**Flow is not a productivity app.**  
**Flow is a ritual for consciousness.**  
**Flow is a boundary you place around presence.**

It creates space for **one thing fully**, then gently returns you to awareness with Zenta.

Following the Unix way: simple, focused, composable.  
Following the Zenta way: mindful, present, untracked.

---

*One session. One focus. One breath at a time.*

