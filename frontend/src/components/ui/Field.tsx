import React from "react";
import { Field as BaseField } from "@base-ui-components/react/field";
import { cn } from "../../utils/cn";

interface FieldRootProps {
  children: React.ReactNode;
  className?: string;
}

const FieldRoot = ({ children, className }: FieldRootProps) => {
  return (
    <BaseField.Root className={cn("space-y-2", className)}>
      {children}
    </BaseField.Root>
  );
};

interface FieldLabelProps {
  children: React.ReactNode;
  className?: string;
}

const FieldLabel = ({ children, className }: FieldLabelProps) => {
  return (
    <BaseField.Label
      className={cn("block text-sm font-medium text-gray-700", className)}
    >
      {children}
    </BaseField.Label>
  );
};

interface FieldControlProps {
  render: React.ReactElement;
}

const FieldControl = ({ render }: FieldControlProps) => {
  return <BaseField.Control render={render} />;
};

interface FieldErrorProps {
  children: React.ReactNode;
  className?: string;
}

const FieldError = ({ children, className }: FieldErrorProps) => {
  return (
    <BaseField.Error className={cn("text-sm text-red-600", className)}>
      {children}
    </BaseField.Error>
  );
};

interface FieldDescriptionProps {
  children: React.ReactNode;
  className?: string;
}

const FieldDescription = ({ children, className }: FieldDescriptionProps) => {
  return (
    <BaseField.Description className={cn("text-sm text-gray-500", className)}>
      {children}
    </BaseField.Description>
  );
};

const Field = {
  Root: FieldRoot,
  Label: FieldLabel,
  Control: FieldControl,
  Error: FieldError,
  Description: FieldDescription,
};

export {
  Field,
  type FieldRootProps,
  type FieldLabelProps,
  type FieldControlProps,
  type FieldErrorProps,
  type FieldDescriptionProps,
};
