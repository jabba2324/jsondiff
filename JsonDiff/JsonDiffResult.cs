namespace JsonDiff
{
    public class JsonDiffResult
    {
        public JsonDiffResult(bool passed, string propertyName, string description)
        {
            Passed = passed;
            PropertyName = propertyName;
            Description = description;
        }

        public bool Passed { get; }

        public bool Failed => !Passed;

        public string PropertyName { get; }

        public string Description { get; }
    }
}
