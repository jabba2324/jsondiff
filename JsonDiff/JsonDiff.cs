using System.Collections.Generic;
using System.Linq;

namespace JsonDiff
{
    public static class JsonDiff
    {
        public static IEnumerable<JsonDiffResult> ComparePropertyNames(Dictionary<string, object> left,
            Dictionary<string, object> right)
        {
            if (left == null)
            {
                return new List<JsonDiffResult>
                {
                    new JsonDiffResult
                    (
                         false,
                         "Request Failed",
                          "Managed body was empty"
                    )
                };
            }

            if (right == null)
            {
                return new List<JsonDiffResult>
                {
                    new JsonDiffResult
                    (
                         false,
                         "Request Failed",
                         "Managed2 body was empty"
                    )
                };
            }

            ICollection<string> leftObjectMembers = left.Keys;
            ICollection<string> rightObjectMembers = right.Keys;

            var passed = leftObjectMembers.Intersect(rightObjectMembers).Select(p => new JsonDiffResult
            (
                 true,
                 p,
                 $"Both Managed & Managed2 have property '{p}'"
            ));

            var leftFailed = leftObjectMembers.Except(rightObjectMembers).Select(p => new JsonDiffResult
            (
                false,
                p,
                $"Error: Managed has property '{p}' but Managed2 does not"
            ));

            var rightFailed = rightObjectMembers.Except(leftObjectMembers).Select(p => new JsonDiffResult
            (
                false,
                p,
                $"Error: Managed2 has property '{p}' but Managed does not"
            ));

            return new List<JsonDiffResult>().Concat(passed).Concat(leftFailed).Concat(rightFailed);
        }

        public static IEnumerable<JsonDiffResult> ComparePropertyValues(IEnumerable<string> propertyNames,
            Dictionary<string, object> left, Dictionary<string, object> right)
        {
            return propertyNames.Select(n =>
            {
                var leftValue = left.Single(m => m.Key == n).Value;
                var rightValue = right.Single(m => m.Key == n).Value;

                var result = CompareForNullValues(leftValue, rightValue, n).ToList();

                if (result.Any())
                {
                    return result;
                }

                switch (leftValue)
                {
                    //Array
                    case List<object> leftArray
                        when rightValue is List<object> rightArray:
                    {
                        if (leftArray.Count == rightArray.Count)
                        { 
                            return Enumerable.Range(0, leftArray.Count)
                                .Select(c => CompareObject((Dictionary<string, object>) leftArray[c],
                                    (Dictionary<string, object>) rightArray[c])).SelectMany(r => r);
                        }

                        return new List<JsonDiffResult>
                        {
                            new JsonDiffResult
                            (
                                  false,
                                  $"Array {n} has {leftArray.Count} items in Managed and {rightArray.Count} items in Managed2",
                                  n
                            )
                        };
                    }

                    //Object
                    case Dictionary<string, object> leftObject
                        when rightValue is Dictionary<string, object> rightObject:
                    {
                        return CompareObject(leftObject, rightObject);
                    }

                    default:
                    {
                        return new List<JsonDiffResult>
                        {
                            ComparePrimitive(n,leftValue,rightValue)
                        };
                    }
                }
            }).SelectMany(r => r);
        }

        private static IEnumerable<JsonDiffResult> CompareForNullValues(object left, object right, string propertyName)
        {
            if (left == null && right != null)
            {
                return new List<JsonDiffResult>
                {
                    new JsonDiffResult
                    (
                         false,
                         propertyName,
                     $"Error: Property '{propertyName}' has value 'null' in Managed and '{right}' in Managed2"
                    )
                }; 
            }

            if (left != null && right == null)
            {
                    return new List<JsonDiffResult>
                    {
                        new JsonDiffResult
                        (
                            false,
                            propertyName,
                    $"Error: Property '{propertyName}' has value '{left}' in Managed and 'null' in Managed2"
                        )
                    };
            }

            if (left == null & right == null)
            {
                return new List<JsonDiffResult>
                {
                    new JsonDiffResult
                    (
                        true,
                        propertyName,
                $"Property '{propertyName}' has value 'null' for both Managed and Managed2"
                    )
                };
            }

            return new List<JsonDiffResult>();
        }

        private static JsonDiffResult ComparePrimitive(string propertyName, object left, object right)
        {
            if (left.Equals(right))
            {
                return new JsonDiffResult
                (
                    true,
                      propertyName,
                      $"Property '{propertyName}' has value '{left}' for both Managed and Managed2"
                );
            }

            return new JsonDiffResult
            (
               false,
                 propertyName,

                    $"Error: Property '{propertyName}' has value '{left}' in Managed and '{right}' in Managed2"
           );
        }

        private static IEnumerable<JsonDiffResult> CompareObject(Dictionary<string, object> left, Dictionary<string, object> right)
        {
            var nameResults = ComparePropertyNames(left, right).ToList();

            var matchingNames = nameResults.Where(t => t.Passed).Select(r => r.PropertyName).ToList();

            var valueResults =
                ComparePropertyValues(matchingNames, left, right);

            return nameResults.Concat(valueResults);
        }
    }
}

