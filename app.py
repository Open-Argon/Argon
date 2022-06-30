import math

def my_sin(theta):
    theta = math.fmod(theta + math.pi, 2 * math.pi) - math.pi
    print(theta)
    result = 0
    termsign = 1
    power = 1

    for i in range(10):
        result += (math.pow(theta, power) / math.factorial(power)) * termsign
        termsign *= -1
        power += 2
    return result

print(my_sin(101))