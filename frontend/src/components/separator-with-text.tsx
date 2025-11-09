import { Separator } from "@/components/ui/separator";

export function SeparatorWithText({ text }: { text: string }) {
  return (
    <div className="relative flex items-center justify-center py-4">
      <Separator className="absolute inset-x-0" />
      <span className="relative z-10 bg-background px-4 text-sm text-muted-foreground">
        {text}
      </span>
    </div>
  );
}
