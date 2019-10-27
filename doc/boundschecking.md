# Bounds Checking Algorithm

## Normalized form

The bounds checking algorithm first normalizes all comparison eqations to this form

`a_0 * x_0 + a_1 * x_1 + ... + a_n * x_n >= 0`

where `a_n` are integers and `x_n` are the product of one or more variables. Here are a few examples of equations converted to normalized form.

|Basic form|Normalized Form|
|----------|---------------|
| a >= b   | a - b >-= 0   |
| a > b    | a - b - 1 >= 0 |
| a == b   | a - b >= 0 && -a + b >= 0 |
| a != b   | a - b - 1 >= 0 \|\| -a + b - 1 >= 0 |
| a * (b + c) >= 0 | a*b + a*c >= 0 |
| a + a > 0 | 2*a - 1 >= 0 |

Using a normalized form ensures that multiple possible ways of writing the same equation end up looking the same to the compiler, such as 
`a + b == c` vs `a == b - c`. Both those forms are normalized to the same thing `a + b - c >= 0 && -a - b + c >= 0`.

## Proving X

Each known polynomial, `a_0 * x_0 + a_1 * x_1 + ... + a_n * x_n >= 0`, can be rewritten as `P_n >= 0`.

Knowing that `P_0 >= 0 && P_1 >= 0` implies that `P_0 + P_1 >= 0`. This also means that given `P_n >= 0` implies that `2 * P_n >= 0`.
Actually, it means that `m * P_n >= 0` for any `m >= 0`.

So given all known polynomials `P_0 >= 0 && P_1 >= 0 && ... && P_n >= 0` you know `X >= 0` if `X = m_0 * P_0 + m_1 * P_1 + ... + m_n * P_n` where `m_0..n >= 0`

As an example. Given `a > b` and `b > c` prove `a > c`.
The first step is to normalize the equations.

Known

`P_0 = 1 >= 0` (universally true)

`P_1 = a - b - 1 >= 0`

`P_2 = b - c - 1 >= 0`

Prove

`X = a - c - 1 >= 0`

```
X = 1 * P_0 + * P_1 + 1 * P_2
  = 1 * 1 + 1 * (a - b - 1) + 1 * (b - c - 1)
  = 1 + a - b - 1 + b - c - 1
  = a - c - 1
```
Since `X` can be expressed in terms of known polynomials
`X` is also true

Expanding the polynomials like this

```
P_0 =  1*1 + 0*a + 0*b + 0*c >= 0
P_1 = -1*1 + 1*a - 1*b + 0*c >= 0
P_2 = -1*1 + 0*a + 1*b - 1*c >= 0
```

Allows us to rewrite these as a matrix

```
    |  1  0  0  0 |
P = | -1  1 -1  0 |
    | -1  0  1 -1 |
```

And `X` can be expressed with a matrix multiply

```
            |  1  0  0  0 |
| 1 1 1 | * | -1  1 -1  0 | = | -1  1  0 -1 |
            | -1  0  1 -1 | 
```
or
```
M * P = X

where 
M = | 1 1 1 |
X = | -1 1 0 -1 |
```
Using this form known constraints can be written as the matrix `P` and the equation you are trying to prove to be true can be written as `X`.

`X` is true if there exists

```
M * P = X
```
Where all values in `M` are non negative
