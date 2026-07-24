# AGENTS — Terminal UI

- Keep `Update` and `View` deterministic; all I/O runs in Bubble Tea commands.
- Recalculate equal-cell layouts from terminal size and card count.
- Mouse clicks and Enter acknowledge changed cards immediately.
- Animate only unacknowledged cards on one shared low-frequency tick.
- Preserve keyboard access for every mouse action and remain readable without color.
