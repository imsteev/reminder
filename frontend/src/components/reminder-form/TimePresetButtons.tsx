import React from "react";
import { Button } from "../ui";
import { TIME_PRESETS } from "../../constants";

interface Props {
  onSetTimeIn: (minutes: number) => void;
  onSetTimeTomorrow: (hour: number, minute: number) => void;
}

export default function TimePresetButtons({ onSetTimeIn, onSetTimeTomorrow }: Props) {
  return (
    <div className="flex gap-2 w-full">
      <Button
        type="button"
        variant="blue"
        size="sm"
        onClick={() => onSetTimeIn(TIME_PRESETS.FIVE_MINUTES)}
        className="rounded-full flex-1"
      >
        5m
      </Button>
      <Button
        type="button"
        variant="blue"
        size="sm"
        onClick={() => onSetTimeIn(TIME_PRESETS.FIFTEEN_MINUTES)}
        className="rounded-full flex-1"
      >
        15m
      </Button>
      <Button
        type="button"
        variant="blue"
        size="sm"
        onClick={() => onSetTimeIn(TIME_PRESETS.ONE_HOUR)}
        className="rounded-full flex-1"
      >
        1h
      </Button>
      <Button
        type="button"
        variant="blue"
        size="sm"
        onClick={() =>
          onSetTimeTomorrow(
            TIME_PRESETS.TOMORROW_9AM.hour,
            TIME_PRESETS.TOMORROW_9AM.minute
          )
        }
        className="rounded-full flex-1"
      >
        Tmr 9am
      </Button>
    </div>
  );
}