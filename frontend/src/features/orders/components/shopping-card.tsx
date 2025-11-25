import {
  AvatarGroup,
  AvatarGroupTooltip,
} from "@/components/animate-ui/components/animate/avatar-group";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Item,
  ItemContent,
  ItemDescription,
  ItemMedia,
  ItemTitle,
} from "@/components/ui/item";
import { AvatarImage } from "@radix-ui/react-avatar";

type Props = {};

const ShoppingCard = (props: Props) => {
  return (
    <div className="flex flex-col gap-3">
      <h3 className="text-2xl font-bold">Shopping List</h3>
      <Card>
        <CardHeader>
          <CardTitle className="flex justify-between text-xl">
            <span>Total items: 3</span>
            <span>
              {Intl.NumberFormat("vi-VN", {
                style: "currency",
                currency: "VND",
              }).format(100000)}
            </span>
          </CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col gap-3 p-4">
          {[1, 2].map(() => (
            <Item variant={"outline"}>
              <ItemMedia>
                <AvatarGroup>
                  <Avatar key={0}>
                    <AvatarImage src={""} />
                    <AvatarFallback>BY</AvatarFallback>
                    <AvatarGroupTooltip>Buyer</AvatarGroupTooltip>
                  </Avatar>

                  <Avatar key={1}>
                    <AvatarImage src={""} />
                    <AvatarFallback>BY</AvatarFallback>
                    <AvatarGroupTooltip>Buyer2</AvatarGroupTooltip>
                  </Avatar>
                </AvatarGroup>
              </ItemMedia>
              <ItemContent>
                <ItemTitle>Test</ItemTitle>
                <ItemDescription>Quantity: 1</ItemDescription>
              </ItemContent>
            </Item>
          ))}
        </CardContent>

        {/* <CardFooter className="flex justify-between">
          <div className="flex items-center gap-2">
            <Avatar className="size-8 sm:size-10">
              <AvatarImage src={""} />
              <AvatarFallback>BY</AvatarFallback>
            </Avatar>
            <h3 className="text-base sm:text-lg leading-tight text-primary">
              Buyer
            </h3>
          </div>

          <QRPayDialog />
        </CardFooter> */}
      </Card>
    </div>
  );
};

export default ShoppingCard;
