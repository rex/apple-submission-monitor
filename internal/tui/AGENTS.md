# AGENTS — Terminal UI

- Keep `Update` and `View` deterministic; all I/O runs in Bubble Tea commands.
- Recalculate equal-cell layouts from terminal size and card count.
- Mouse clicks and Enter acknowledge changed cards immediately.
- Animate cached hero gradients on one shared low-frequency tick; reserve
  border/background flashing for unacknowledged status changes.
- Preserve keyboard access for every mouse action and remain readable without color.
- Treat approved outcomes as a dedicated victory state with a giant outcome
  hero; never leave a completed card visually frozen in its prior status.
