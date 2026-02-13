/**
 * # Example Field Component
 *
 * Örnek custom field component'i.
 * Plugin sisteminde custom field nasıl oluşturulur gösterir.
 *
 * ## Kullanım
 * Backend'de field tanımı:
 * ```go
 * fields.Custom("Example Field", "example_field").
 *     Type("example-field").
 *     Label("Example Field")
 * ```
 */

import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

interface ExampleFieldProps {
  field: {
    key: string;
    label?: string;
    placeholder?: string;
    helpText?: string;
    required?: boolean;
  };
  value: string;
  onChange: (value: string) => void;
  error?: string;
  disabled?: boolean;
}

/**
 * ExampleField Component
 *
 * Custom field component örneği.
 * Standart input field'ı card içinde gösterir.
 */
export function ExampleField({
  field,
  value,
  onChange,
  error,
  disabled,
}: ExampleFieldProps) {
  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle className="text-sm font-medium">
          {field.label || field.key}
          {field.required && <span className="text-destructive ml-1">*</span>}
        </CardTitle>
        {field.helpText && (
          <CardDescription className="text-xs">
            {field.helpText}
          </CardDescription>
        )}
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          <Label htmlFor={field.key} className="sr-only">
            {field.label || field.key}
          </Label>
          <Input
            id={field.key}
            type="text"
            value={value || ''}
            onChange={(e) => onChange(e.target.value)}
            placeholder={field.placeholder || `Enter ${field.label || field.key}`}
            disabled={disabled}
            className={error ? 'border-destructive' : ''}
          />
          {error && (
            <p className="text-xs text-destructive mt-1">{error}</p>
          )}
          <p className="text-xs text-muted-foreground">
            This is an example custom field from ExamplePlugin
          </p>
        </div>
      </CardContent>
    </Card>
  );
}
