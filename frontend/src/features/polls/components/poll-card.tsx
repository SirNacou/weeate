import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
} from "@/components/ui/card";
import PollStrategyBadge from "./poll-strategy-badge";
import CloseTimer from "@/components/close-timer";
import { add } from "date-fns";
import { RadioGroup } from "@/components/ui/radio-group";
import { useForm } from "@tanstack/react-form";
import { z } from "zod";
import { useMemo, useState } from "react";
import { FieldGroup } from "@/components/ui/field";
import { Button } from "@/components/ui/button";
import PollOption, { Option } from "./poll-option";
import { Spinner } from "@/components/ui/spinner";
import { toast } from "sonner";

type Props = {
  buyerName: string;
  avatarUrl: string;
  strategy: "ORDER_CONSENSUS_ITEM" | "ORDER_PERSONAL_CHOICE";
  options: Option[];
  initialSelectedOption?: string;
};

const PollCard = ({
  buyerName,
  avatarUrl,
  strategy,
  options,
  initialSelectedOption,
}: Props) => {
  const [now] = useState(Date.now());

  const validator = useMemo(() => {
    return z.object({
      selectedOption: z
        .string()
        .nonempty({ message: "Please select an option" })
        .refine((val) => options.some((option) => option.id === val), {
          message: "Invalid option selected",
        }),
    });
  }, [options]);

  const form = useForm({
    defaultValues: {
      selectedOption: initialSelectedOption || "",
    },
    validators: {
      onChange: validator,
    },
    listeners: {
      onChange(props) {
        console.log("Form changed:", props.fieldApi.state.value);
      },
    },
    onSubmit: async ({ value }) => {
      return new Promise((resolve) => setTimeout(resolve, 2000)).then(() => {
        toast.success("Vote submitted successfully!");
        console.log("Submitted value:", value);
      });
    },
  });

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();
        e.stopPropagation();
        form.handleSubmit();
      }}
    >
      <Card>
        <CardHeader>
          <CardTitle className="flex flex-col sm:flex-row text-lg">
            <div className="flex items-center gap-2">
              <Avatar className="size-8 sm:size-10">
                <AvatarImage src={avatarUrl} />
                <AvatarFallback>
                  {buyerName.substring(0, 2).toUpperCase()}
                </AvatarFallback>
              </Avatar>
              <h3 className="text-base sm:text-lg leading-tight">
                <span className="text-primary font-semibold">{buyerName}</span>{" "}
                is buying breakfast
              </h3>
            </div>
          </CardTitle>
          <CardDescription className="mt-1 flex items-center justify-between text-sm sm:text-base">
            <CloseTimer
              className="text-sm sm:text-base shrink-0"
              closesAt={add(now, { minutes: 1 })}
            />

            <PollStrategyBadge strategy={strategy} />
          </CardDescription>
        </CardHeader>
        <CardContent className="p-4 sm:p-6 pt-0 sm:pt-0">
          <FieldGroup>
            <form.Field
              name="selectedOption"
              children={(field) => {
                return (
                  <RadioGroup
                    className={"grid grid-cols-1 lg:grid-cols-2 gap-3"}
                    value={field.state.value}
                    onValueChange={field.handleChange}
                  >
                    {options.map((option) => (
                      <PollOption
                        key={option.id}
                        option={option}
                        isSelected={option.id === field.state.value}
                        onSelect={(id) => field.handleChange(id)}
                        disabled={false}
                      />
                    ))}
                  </RadioGroup>
                );
              }}
            />
          </FieldGroup>
        </CardContent>
        <CardFooter>
          <form.Subscribe
            selector={(state) => [
              state.isSubmitting,
              state.canSubmit,
              state.isPristine,
            ]}
            children={([isSubmitting, canSubmit, isPristine]) => (
              <Button
                type="submit"
                className="w-full"
                disabled={!canSubmit || isPristine}
              >
                {isSubmitting ?
                  <>
                    <Spinner />
                    Submitting...
                  </>
                : "Submit"}
              </Button>
            )}
          />
        </CardFooter>
      </Card>
    </form>
  );
};

export default PollCard;
