package products

type GetProductInput struct {
    ID string
}

type GetProductOutput struct {
    ID    string
    Name  string
    Price float64
}

func (c *component) GetProduct(in GetProductInput) (GetProductOutput, error) {
    p, err := c.getUC.Execute(in.ID) 
    if err != nil {
        return GetProductOutput{}, err
    }

    return GetProductOutput{
        ID:    p.ID,
        Name:  p.Name,
        Price: p.Price,
    }, nil
}