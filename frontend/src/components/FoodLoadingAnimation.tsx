import { motion } from "motion/react";
import { Icon } from "@iconify/react";

export default function FoodLoadingAnimation() {
  const foods = [
    { icon: "noto:steaming-bowl", delay: 0 }, // Phở
    { icon: "noto:baguette-bread", delay: 0.2 }, // Bánh mì
    { icon: "noto:pot-of-food", delay: 0.4 }, // Hot pot
    { icon: "noto:leafy-green", delay: 0.6 }, // Salad/Gỏi
    { icon: "noto:cooked-rice", delay: 0.8 }, // Cơm
  ];

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-linear-to-br from-orange-50 to-red-50">
      <div className="flex flex-col items-center gap-8">
        {/* Animated food icons */}
        <div className="flex items-center gap-4">
          {foods.map((food, index) => (
            <motion.div
              key={index}
              className="text-5xl"
              animate={{
                y: [0, -20, 0],
              }}
              transition={{
                duration: 0.8,
                repeat: Infinity,
                ease: "easeInOut",
                delay: food.delay,
              }}
            >
              <Icon icon={food.icon} width={48} height={48} />
            </motion.div>
          ))}
        </div>

        {/* Loading text */}
        <motion.p
          className="text-xl font-semibold text-gray-700"
          animate={{
            opacity: [0.5, 1, 0.5],
          }}
          transition={{
            duration: 2,
            repeat: Infinity,
            ease: "easeInOut",
          }}
        >
          Loading...
        </motion.p>

        {/* Progress bar */}
        <div className="w-48 h-2 bg-gray-200 rounded-full overflow-hidden">
          <motion.div
            className="h-full bg-orange-500"
            animate={{
              x: ["-100%", "100%"],
            }}
            transition={{
              duration: 1.5,
              repeat: Infinity,
              ease: "easeInOut",
            }}
            style={{ width: "50%" }}
          />
        </div>
      </div>
    </div>
  );
}
