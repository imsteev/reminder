import React, { forwardRef, useEffect, useState } from "react";
import { Input, type InputProps } from "./Input";

interface Props extends Omit<InputProps, "type" | "value" | "onChange"> {
  value?: string;
  onChange?: (value: string) => void;
}

const PhoneInput = forwardRef<HTMLInputElement, Props>(
  ({ value = "", onChange, onBlur, ...props }, ref) => {
    const [displayValue, setDisplayValue] = useState(() =>
      formatPhoneNumber(value)
    );

    useEffect(() => {
      setDisplayValue(formatPhoneNumber(value));
    }, [value]);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      const newValue = e.target.value;
      const cleaned = newValue.replace(/\D/g, "");

      if (cleaned.length <= 10) {
        const formatted = formatPhoneNumber(newValue);
        setDisplayValue(formatted);
        onChange?.(getUnformattedValue(formatted));
      }
    };

    const handleBlur = (e: React.FocusEvent<HTMLInputElement>) => {
      const cleaned = getUnformattedValue(displayValue);
      if (cleaned.length > 0 && cleaned.length < 10) {
        console.warn("Phone number should be 10 digits");
      }
      onBlur?.(e);
    };

    return (
      <Input
        {...props}
        ref={ref}
        type="tel"
        value={displayValue}
        onChange={handleChange}
        onBlur={handleBlur}
        placeholder="(555) 123-4567"
        maxLength={14}
      />
    );
  }
);

export const formatPhoneNumber = (value: string): string => {
  const cleaned = value.replace(/\D/g, "");

  if (cleaned.length === 0) {
    return "";
  }
  if (cleaned.length <= 3) {
    return `(${cleaned}`;
  }
  if (cleaned.length <= 6) {
    return `(${cleaned.slice(0, 3)}) ${cleaned.slice(3)}`;
  }

  const parts = [
    cleaned.slice(0, 3),
    cleaned.slice(3, 6),
    cleaned.slice(6, 10),
  ];
  return `(${parts[0]}) ${parts[1]}-${parts[2]}`;
};

const getUnformattedValue = (value: string): string => {
  return value.replace(/\D/g, "");
};

PhoneInput.displayName = "PhoneInput";

export { PhoneInput, type Props as PhoneInputProps };
