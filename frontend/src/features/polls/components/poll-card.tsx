import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "@/components/ui/card";
import PollStrategyBadge from "./poll-strategy-badge";
import CloseTimer from "@/components/close-timer";
import { add } from "date-fns";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { cn } from "@/lib/utils";
import { useForm } from "@tanstack/react-form";
import zod from 'zod'
import { form } from "motion/react-client"

type Props = {
  buyerName: string;
  avatarUrl: string;
  strategy: "ORDER_CONSENSUS_ITEM" | "ORDER_PERSONAL_CHOICE";
  options: string[];
};

const formValidators = zod.object({
  selectedOption: zod.string().nonempty("Please select an option"),
})

const PollCard = ({ buyerName, avatarUrl, strategy, options }: Props) => {
  const gridColsClass = options.length > 6 ? "grid-cols-3" : "grid-cols-2";

  const form = useForm({
    defaultValues: {
      selectedOption: null,
    },
    validators:
    onSubmit: (values) => {
      console.log(values);
    },
  });

  return (
    <Card>
      <form
        onSubmit={(e) => {
          e.preventDefault();
          e.stopPropagation();
          form
        }}
      ></form>
      <CardHeader>
        <CardTitle className="flex gap-2 items-center text-lg">
          <Avatar className="size-10">
            <AvatarImage src={avatarUrl} />
            <AvatarFallback>
              {buyerName.substring(0, 2).toUpperCase()}
            </AvatarFallback>
          </Avatar>
          <h3 className="text-lg leading-none">
            <span className="text-primary">{buyerName}</span> is buying
            breakfast
          </h3>

          <CloseTimer
            className="ms-auto text-sm"
            closesAt={add(Date.now(), { minutes: 10 })}
          />
        </CardTitle>
        <CardDescription>
          <PollStrategyBadge strategy={strategy} />
        </CardDescription>
      </CardHeader>
      <CardContent>
        <RadioGroup>
          <div className={cn("grid gap-4", gridColsClass)}>
            {options.map((option, index) => (
              <div className="cursor-pointer select-none">
                <RadioGroupItem
                  id={option}
                  key={index}
                  value={option}
                  className="mb-2 last:mb-0"
                />
                <Label htmlFor={option}>{option}</Label>
              </div>
            ))}
          </div>
        </RadioGroup>
        {/* Replace with actual poll list */}
        <p>No active polls at the moment.</p>
      </CardContent>
    </Card>
  );
};

export default PollCard;
