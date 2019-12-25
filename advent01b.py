def fuelweight(volume):
    weight = (volume // 3) - 2
    if weight <= 0:
        return 0
    return weight + fuelweight(weight)

if __name__ == "__main__":
    print sum(fuelweight(int(line)) for line in open("advent01a.txt"))
