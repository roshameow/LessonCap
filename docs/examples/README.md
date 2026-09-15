# Synthetic capture exercise

Open `slides.html` in a browser, select the slide rectangle in LessonCap and start monitoring. Advance each page manually, waiting several seconds between changes. Stop capture and export. Inspect the resulting page count, order, crop and readability; record duplicates or omissions.

`sample-slides.pdf` is a reference document generated with ReportLab from synthetic content. It illustrates page layout and can also be opened as a source for manual capture. It is **not** an output recorded from the native LessonCap app, nor evidence of capture accuracy. Screen resolution and the selected region determine the actual export dimensions.

Regenerate the reference PDF with:

```bash
python -m pip install reportlab
python scripts/generate_example.py
```

For a short product recording, show the selected area, page changes, increasing capture count, then the exported PDF. Record the actual app and include the OS and app version in the recording description.
